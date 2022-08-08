package datatunnel

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/zhiting-tech/smartassistant/pkg/datatunnel/v2/control"
	"github.com/zhiting-tech/smartassistant/pkg/datatunnel/v2/proto"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	DefaultCallTimeout = 10 * time.Second

	DefaultGrpcIdleTime       = 10 * time.Second
	DefaultGrpcPingAckTimeout = 20 * time.Second
)

type ProxyService struct {
	ServiceName string
	ServicePort int
	ServiceHost string
}

type ProxyControlClientHook struct {
	OnConnected func(c *ProxyControlClient, pcsc *control.ProxyControlStreamContext) error
}

type ProxyControlClientOptionFn func(*ProxyControlClient)

func (fn ProxyControlClientOptionFn) apply(c *ProxyControlClient) {
	fn(c)
}

func WithProxyAddr(addr string) ProxyControlClientOptionFn {
	return func(pcc *ProxyControlClient) {
		pcc.proxyAddr = addr
	}
}

func WithClientHook(hook ProxyControlClientHook) ProxyControlClientOptionFn {
	return func(pcc *ProxyControlClient) {
		pcc.hook = hook
	}
}

func WithLogger(logger control.Logger) ProxyControlClientOptionFn {
	return func(pcc *ProxyControlClient) {
		pcc.logger = logger
	}
}

type ProxyControlClient struct {
	base        *control.ControlBase
	version     int32
	hook        ProxyControlClientHook
	logger      control.Logger
	servicesMap map[string]*ProxyService
	proxyAddr   string
	currentPCSC atomic.Value
}

func NewProxyControlClient(opts ...ProxyControlClientOptionFn) *ProxyControlClient {
	// 旧版本SA编译出的protobuf是messageV1版本，此时由于旧版本coder不支持messageV1，因此使用的是json序列化
	// 新版本SA编译出的protobuf是messageV2版本，因此使用protobuf序列化
	// 开启后会导致使用protobuf序列化而无法兼容旧版本，因此暂时关闭支持protobuf结构体
	b := &control.DefaultBinaryCoder{EnableProtobuf: false}
	client := &ProxyControlClient{
		logger:      logrus.New(),
		version:     clientVersion,
		servicesMap: map[string]*ProxyService{},
	}

	for _, opt := range opts {
		opt.apply(client)
	}

	client.base = control.NewControlBase(
		control.WithCoder(b),
		control.WithLogger(client.logger),
	)

	client.init()
	return client
}

func (c *ProxyControlClient) init() {
	c.base.RegisterClientMethod(methodAuthenticate, c.version, proto.ProxyControlStreamMsgType_REQUEST, c.Authenticate)
	c.base.RegisterClientMethod(methodRegisterService, c.version, proto.ProxyControlStreamMsgType_REQUEST, c.RegisterService)
	c.base.RegisterClientMethod(methodAllowTempConnCert, c.version, proto.ProxyControlStreamMsgType_REQUEST, c.AllowTempConnCert)
	c.base.RegisterRPC(notifyNewConnection, c.version, proto.ProxyControlStreamMsgType_NOTIFY, c.NewConnection)
}

func (c *ProxyControlClient) RegisterProxyServices(services ...ProxyService) {
	for _, service := range services {
		c.servicesMap[service.ServiceName] = &ProxyService{
			ServiceName: service.ServiceName,
			ServicePort: service.ServicePort,
			ServiceHost: service.ServiceHost,
		}
	}
}

func (c *ProxyControlClient) Run(ctx context.Context, target string) {
	var (
		pcsc   *control.ProxyControlStreamContext
		client proto.ProxyController_ControlStreamClient
		stop   chan struct{} = make(chan struct{})
		err    error
	)
	c.logger.Debug("run proxy control client")
	for {
		if pcsc, client, err = c.DialWithContext(ctx, target); err != nil {
			c.logger.Debugf("Dial error %v", err)
		} else {
			go func() {
				var err error
				if err = c.onConnected(pcsc); err != nil {
					pcsc.Close()

					c.logger.Debugf("OnConnected error %v", err)

					if cerr, ok := err.(*control.ControlError); ok {
						if cerr.GetCode() == InvalidSAIDOrSAKey {
							close(stop)

							c.logger.Infof("InvalidSAIDOrSAKey, now break loop")
						}
					}
				} else {
					c.logger.Debugf("OnConnected success")
					c.currentPCSC.Store(pcsc)
				}

			}()
			c.HandleMsg(pcsc, client)
		}

		select {
		case <-ctx.Done():
			return
		case <-stop:
			c.logger.Info("invalid said or sakey, so stop proxy")
			return
		case <-time.After(10 * time.Second):
			c.logger.Infof("retry connect to %s", target)
		}
	}
}

func (c *ProxyControlClient) GetCurrentPCSC() (pcsc *control.ProxyControlStreamContext, err error) {
	value := c.currentPCSC.Load()
	if value == nil {
		err = errors.New("proxy control client not connected")
		return
	}

	pcsc = value.(*control.ProxyControlStreamContext)
	return
}

func (c *ProxyControlClient) DialWithContext(ctx context.Context, target string) (
	pcsc *control.ProxyControlStreamContext,
	client proto.ProxyController_ControlStreamClient,
	err error,
) {
	var (
		conn *grpc.ClientConn
	)

	// 空闲时间10秒则发送ping,发送ping后20秒内收不到ack视为断开连接
	kacp := keepalive.ClientParameters{
		Time:                DefaultGrpcIdleTime,
		Timeout:             DefaultGrpcPingAckTimeout,
		PermitWithoutStream: true,
	}

	if conn, err = grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp)); err != nil {
		return
	}

	if client, err = proto.NewProxyControllerClient(conn).ControlStream(ctx); err != nil {
		return
	}

	pcsc = control.NewProxyControlStreamContextWithContext(ctx, client)
	return
}

func (c *ProxyControlClient) onConnected(pcsc *control.ProxyControlStreamContext) (err error) {
	if c.hook.OnConnected != nil {
		if err = c.hook.OnConnected(c, pcsc); err != nil {
			return
		}
	}

	err = c.RegisterAllServices(pcsc)

	return
}

func (c *ProxyControlClient) HandleMsg(pcsc *control.ProxyControlStreamContext, client proto.ProxyController_ControlStreamClient) {
	var (
		msg *proto.ProxyControlStreamMsg
		err error
	)
	for {
		if msg, err = client.Recv(); err != nil {
			pcsc.Close()
			return
		}

		c.base.HandleProxyControlStreamMsg(pcsc, msg)
	}
}

func (c *ProxyControlClient) Authenticate(pcsc *control.ProxyControlStreamContext, req *proto.AuthenticateRequest) (err error) {
	var (
		caller *control.RemoteCaller
	)

	if caller, err = c.base.NewRemoteCaller(methodAuthenticate, c.version); err != nil {
		return
	}

	_, err = caller.Call(pcsc, req)

	return
}

func (c *ProxyControlClient) RegisterAllServices(pcsc *control.ProxyControlStreamContext) (err error) {
	var (
		req *proto.RegisterServiceRequest = &proto.RegisterServiceRequest{
			Services: []*proto.RegisterServiceValue{},
		}
	)

	for _, service := range c.servicesMap {
		req.Services = append(req.Services, &proto.RegisterServiceValue{
			ServiceName: service.ServiceName,
			ServicePort: int32(service.ServicePort),
		})
	}

	return c.RegisterService(pcsc, req)
}

func (c *ProxyControlClient) RegisterService(pcsc *control.ProxyControlStreamContext, req *proto.RegisterServiceRequest) (err error) {
	var (
		caller *control.RemoteCaller
	)

	if caller, err = c.base.NewRemoteCaller(methodRegisterService, c.version); err != nil {
		return
	}

	_, err = caller.Call(pcsc, req)

	return
}

func (c *ProxyControlClient) NewConnection(pcsc *control.ProxyControlStreamContext, event *proto.NewConnectionEvent) (err error) {
	service, ok := c.servicesMap[event.ServiceName]
	if !ok {
		err = fmt.Errorf("service %s not found", event.ServiceName)
		return
	}

	addr := c.proxyAddr
	if event.RemoteAddr != "" {
		addr = event.RemoteAddr
	}

	if addr == "" {
		err = fmt.Errorf("server addr is empty")
		return
	}

	c.proxyService(event.Key, addr, fmt.Sprintf("%s:%d", service.ServiceHost, service.ServicePort))
	return
}

func (c *ProxyControlClient) proxyService(key []byte, serverAddr string, serviceAddr string) {
	client := NewProxyClient(key)
	ch := make(chan struct{}, 2)
	var conn1 net.Conn
	go func() {
		var err error
		if conn1, err = client.DialContext(context.TODO(), "tcp", serverAddr); err != nil {
			c.logger.Warnf("dial %s error %v", serverAddr, err)
		}
		c.logger.Debugf("dial %s finish", serverAddr)
		ch <- struct{}{}
	}()

	var conn2 net.Conn
	go func() {
		var err error
		if conn2, err = net.Dial("tcp", serviceAddr); err != nil {
			c.logger.Warnf("dial %s error %v", serviceAddr, err)
		}
		c.logger.Debugf("dial %s finish", serviceAddr)
		ch <- struct{}{}
	}()

	<-ch
	<-ch
	close(ch)

	if conn1 == nil || conn2 == nil {
		if conn1 != nil {
			conn1.Close()
		}

		if conn2 != nil {
			conn2.Close()
		}
		return
	}

	proxy := func(dst, src net.Conn) {
		defer dst.Close()
		defer src.Close()
		io.Copy(dst, src)
	}

	go proxy(conn1, conn2)
	go proxy(conn2, conn1)

}
