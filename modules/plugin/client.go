package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/zhiting-tech/smartassistant/modules/types/status"
	errors2 "github.com/zhiting-tech/smartassistant/pkg/errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/zhiting-tech/smartassistant/modules/entity"
	event2 "github.com/zhiting-tech/smartassistant/pkg/event"
	"github.com/zhiting-tech/smartassistant/pkg/logger"
	"github.com/zhiting-tech/smartassistant/pkg/plugin/sdk/proto/v2"
	"github.com/zhiting-tech/smartassistant/pkg/plugin/sdk/v2"
	"github.com/zhiting-tech/smartassistant/pkg/plugin/sdk/v2/definer"
	"github.com/zhiting-tech/smartassistant/pkg/thingmodel"
)

var NotExistErr = errors.New("plugin not exist")

func NewClient() *client {
	return &client{
		clients: make(map[string]*pluginClient),
	}
}

type client struct {
	mu      sync.Mutex // clients 锁
	clients map[string]*pluginClient
}

// DevicesDiscover 发现设备，并且通过 channel 返回给调用者
func (c *client) DevicesDiscover(ctx context.Context) <-chan DiscoverResponse {
	out := make(chan DiscoverResponse, 10)
	go func() {
		var wg sync.WaitGroup
		for _, cli := range c.clients {
			wg.Add(1)
			go func(cli *pluginClient) {
				defer wg.Done()
				logger.Debug("listening plugin Discovering...")
				cli.DeviceDiscover(ctx, out)
				logger.Debug("plugin listening done")
			}(cli)
		}
		wg.Wait()
		close(out)
	}()
	return out
}

func (c *client) Configs() (configs []Plugin) {

	for _, cli := range c.clients {
		configs = append(configs, cli.PluginConf)
	}
	return
}

func (c *client) Config(pluginID string) (config Plugin) {

	cli, err := c.get(pluginID)
	if err != nil {
		return
	}
	return cli.PluginConf
}

// Connect 连接设备
func (c *client) Connect(ctx context.Context, identify Identify, authParams map[string]interface{}) (das thingmodel.ThingModel, err error) {

	pc, err := c.get(identify.PluginID)
	if err != nil {
		return
	}
	d := pc.Device(identify.IID)
	das, err = d.Connect(ctx, authParams)
	if err != nil {
		return
	}
	if err = d.WaitOnline(ctx); err != nil {
		return thingmodel.ThingModel{}, errors2.New(status.DeviceConnectTimeout)
	}
	for _, ins := range das.Instances {
		dd := pc.Device(ins.IID)
		dd.HealthCheck()
	}
	return
}

func (c *client) Disconnect(ctx context.Context, identify Identify, authParams map[string]interface{}) (err error) {
	cli, err := c.get(identify.PluginID)
	if err != nil {
		return
	}
	return cli.RemoveDevice(ctx, identify.IID, authParams)
}

func (c *client) SetAttributes(ctx context.Context, pluginID string, areaID uint64, setReq sdk.SetRequest) (result []byte, err error) {
	data, _ := json.Marshal(setReq)
	req := proto.SetAttributesReq{
		Data: data,
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	cli, err := c.get(pluginID)
	if err != nil {
		return
	}
	_, err = cli.protoClient.SetAttributes(ctx, &req)
	if err != nil {
		return
	}
	return
}

func (c *client) IsOnline(identify Identify) bool {
	cli, err := c.get(identify.PluginID)
	if err != nil {
		logger.Warningf("plugin %s not found", identify.PluginID)
		return false
	}
	d := cli.Device(identify.IID)
	d.HealthCheck()
	return d.IsOnline()
}

func (c *client) OTA(ctx context.Context, identify Identify, firmwareURL string) (err error) {

	logger.Infof("ota: %s, firmware url: %s", identify.IID, firmwareURL)
	cli, err := c.get(identify.PluginID)
	if err != nil {
		return
	}

	return cli.Device(identify.IID).OTA(ctx, firmwareURL)
}

func (c *client) get(domain string) (*pluginClient, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cli, ok := c.clients[domain]
	if ok {
		return cli, nil
	}
	return nil, errors2.New(status.PluginDomainNotExist)
}

func (c *client) Add(cli *pluginClient) {

	c.mu.Lock()
	defer c.mu.Unlock()
	oldCli, ok := c.clients[cli.pluginID]
	if ok {
		go oldCli.Stop()
	}
	c.clients[cli.pluginID] = cli
	go cli.InitDevices()
	go c.ListenStateChange(cli.pluginID)
}

func (c *client) Remove(pluginID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	cli, ok := c.clients[pluginID]
	if ok {
		delete(c.clients, pluginID)
		go cli.Stop()
	}
	return nil
}

func (c *client) ListenStateChange(pluginID string) {
	for {
		cli, err := c.get(pluginID)
		if err != nil {
			logrus.Errorf("get %s's plugin client err: %s", pluginID, err)
			return
		}
		if err = cli.ListenChange(); err != nil {
			logrus.Errorf("%s listen state change err: %s", pluginID, err)
		}

		logrus.Warningf("%s try listen state change", pluginID)
		time.Sleep(time.Second * 2)
	}
}

func HandleEvent(cli *pluginClient, ev sdk.Event) (err error) {
	switch ev.Type {
	case sdk.AttrChangeEvent:
		var d entity.Device
		var attrEvent definer.AttributeEvent
		_ = json.Unmarshal(ev.Data, &attrEvent)

		d, err = entity.GetPluginDevice(cli.areaID, cli.pluginID, attrEvent.IID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			logger.Errorf("GetPluginDevice error:%s", err.Error())
			return
		}

		em := event2.NewEventMessage(event2.AttributeChange, cli.areaID)
		em.SetDeviceID(d.ID)
		em.SetAttr(attrEvent)
		event2.Notify(em)
	case sdk.ThingModelChangeEvent:
		var tme definer.ThingModelEvent
		_ = json.Unmarshal(ev.Data, &tme)
		em := event2.NewEventMessage(event2.ThingModelChange, cli.areaID)
		em.Param = map[string]interface{}{
			"thing_model": tme.ThingModel,
			"iid":         tme.IID,
			"area_id":     strconv.FormatUint(cli.areaID, 10),
			"plugin_id":   cli.pluginID,
		}
		event2.Notify(em)
	}
	return
}

func ParseAttrsResp(resp *proto.GetInstancesResp) thingmodel.ThingModel {

	var instances []thingmodel.Instance
	_ = json.Unmarshal(resp.Instances, &instances)

	var authParams []thingmodel.AuthParam
	if resp.AuthRequired {
		if err := json.Unmarshal(resp.AuthParams, &authParams); err != nil {
			logrus.Errorf("unmarshal authParams err: %s", err)
		}
	}
	return thingmodel.ThingModel{
		Instances:  instances,
		OTASupport: resp.OtaSupport,

		AuthRequired: resp.AuthRequired,
		IsAuth:       resp.IsAuth,
		AuthParams:   authParams,
	}
}
