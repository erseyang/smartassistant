package registry

import (
	"context"
	"fmt"
	"github.com/zhiting-tech/smartassistant/pkg/plugin/sdk/utils/addr"
	"net"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"

	"github.com/zhiting-tech/smartassistant/pkg/logger"
)

const (
	registerTTL = 10

	etcdURL = "http://0.0.0.0:2379"

	ManagerTarget = "/sa/plugins"
)

func EndpointsKey(service string) string {
	return fmt.Sprintf("%s/%s", ManagerTarget, service)
}

// RegisterService 注册服务
func RegisterService(ctx context.Context, key string, endpoint endpoints.Endpoint) {
	logger.Info("register service:", key, endpoint.Addr)
	cli, err := clientv3.NewFromURL(etcdURL)
	if err != nil {
		logger.Errorf("new etcd client err: %s", err.Error())
		return
	}
	defer cli.Close()
	for {
		if err = register(ctx, cli, key, endpoint); err != nil {
			logger.Errorf("register service err: %s", err.Error())
		}
		time.Sleep(time.Second)
	}
}

func register(ctx context.Context, cli *clientv3.Client, key string, endpoint endpoints.Endpoint) (err error) {
	kl, lease, grant, err := CreateLease(ctx, cli, key, endpoint)
	if err != nil {
		return
	}
	for {
		if _, ok := <-kl; !ok {
			lease.Close()
			time.Sleep(time.Second)
			kl, lease, grant, err = CreateLease(ctx, cli, key, endpoint)
			if err != nil {
				return
			}
		}
		if IsIPChange(endpoint) {
			UpdateEndPoint(ctx, cli, key, &endpoint, grant.ID)
		}
	}
}

// UnregisterService 取消注册服务
func UnregisterService(ctx context.Context, key string) (err error) {
	logger.Info("unregister service:", key)
	cli, err := clientv3.NewFromURL(etcdURL)
	if err != nil {
		return
	}
	em, err := endpoints.NewManager(cli, ManagerTarget)
	if err != nil {
		return
	}

	return em.DeleteEndpoint(ctx, key)
}

func CreateLease(ctx context.Context, cli *clientv3.Client, key string, endpoint endpoints.Endpoint) (
	kl <-chan *clientv3.LeaseKeepAliveResponse, lease clientv3.Lease, grant *clientv3.LeaseGrantResponse, err error ){
	em, err := endpoints.NewManager(cli, ManagerTarget)
	if err != nil {
		return
	}

	lease = clientv3.NewLease(cli)
	grant, err = lease.Grant(ctx, registerTTL)
	if err != nil {
		return
	}
	kl, err = lease.KeepAlive(ctx, grant.ID)
	if err != nil {
		return
	}

	err = em.AddEndpoint(ctx, key, endpoint, clientv3.WithLease(grant.ID))
	if err != nil {
		return
	}
	return
}

func IsIPChange(endpoint endpoints.Endpoint) bool {
	oldAddr, err := net.ResolveTCPAddr("tcp", endpoint.Addr)
	if err != nil {
		logger.Errorf("IsIpChange ResolveTCPAddr err: %s", err)
		return false
	}
	newIP, err := addr.LocalIP()
	if err != nil {
		logger.Errorf("IsIpChange get LocalIP err: %s", err)
		return false
	}
	if newIP == oldAddr.IP.String() {
		return false
	}
	return true
}

func UpdateEndPoint(ctx context.Context, cli *clientv3.Client, key string, endpoint *endpoints.Endpoint, leaseId clientv3.LeaseID) {
	newIP, err := addr.LocalIP()
	if err != nil {
		logger.Errorf("UpdateEndPoint get LocalIP err: %s", err)
		return
	}
	oldAddr, err := net.ResolveTCPAddr("tcp", endpoint.Addr)
	if err != nil {
		logger.Errorf("UpdateEndPoint ResolveTCPAddr err: %s", err)
		return
	}
	ipAddr := net.TCPAddr{
		IP:   net.ParseIP(newIP),
		Port: oldAddr.Port,
	}
	endpoint.Addr = ipAddr.String()
	em, err := endpoints.NewManager(cli, ManagerTarget)
	if err != nil {
		logger.Errorf("UpdateEndPoint NewManager err: %s", err)
		return
	}
	err = em.AddEndpoint(ctx, key, *endpoint, clientv3.WithLease(leaseId))
	if err != nil {
		logger.Errorf("UpdateEndPoint AddEndpoint err: %s", err)
		return
	}
	return
}