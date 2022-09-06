package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"

	"github.com/zhiting-tech/smartassistant/pkg/plugin/sdk/v2/definer"
	"github.com/zhiting-tech/smartassistant/pkg/thingmodel"
)

var (
	DeviceNotExist = errors.New("device not exist")
)

const (
	ThingModelChangeEvent = "thing_model_change"
	AttrChangeEvent       = "attr_change"
)

func newDevice(d Device) *device {
	dd := device{
		Device: d,
	}
	return &dd
}

type device struct {
	Device
	df *definer.Definer

	connected bool
	mutex     sync.Mutex
}

type Manager struct {
	devices         sync.Map // iid: device
	eventChannel    EventChan
	eventChannels   map[EventChan]struct{}
	discoverDevices sync.Map // iid: device
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) GetDiscoverDevice(iid string) *device {
	v, ok := m.discoverDevices.Load(iid)
	if ok {
		return v.(*device)
	}
	return nil
}

func (m *Manager) Init() {
	m.eventChannel = make(EventChan, 10)
	m.eventChannels = make(map[EventChan]struct{})

	// 转发notifyChan消息到所有notifyChans
	go func() {
		for {
			select {
			case n := <-m.eventChannel:
				for ch := range m.eventChannels {
					select {
					case ch <- n:
					default:
					}
				}
			}
		}
	}()
}

func (m *Manager) Subscribe(ch EventChan) {
	m.eventChannels[ch] = struct{}{}
}

func (m *Manager) Unsubscribe(ch EventChan) {
	delete(m.eventChannels, ch)
}

// InitOrUpdateDevice 添加或更新设备, 更新设备将会是以重连的形式更新，原有设备会被回收资源
func (m *Manager) InitOrUpdateDevice(d Device) error {
	if d == nil {
		return errors.New("device is nil")
	}

	v, loaded := m.devices.LoadOrStore(d.Info().IID, newDevice(d))
	if loaded {
		logrus.Debugf("device %s already exist", d.Info().IID)
		oldDevice := v.(*device)
		addrChange := oldDevice.Address() != d.Address()
		if addrChange {
			logrus.Debugf("device %s change address: %s -> %s:",
				d.Info().IID, oldDevice.Address(), d.Address())
		}
		disconnected := oldDevice.connected && !oldDevice.Online(d.Info().IID)
		if disconnected {
			logrus.Debugf("device %s disconnected, reconnecting...", d.Info().IID)
		}
		if addrChange || disconnected {
			m.devices.Store(d.Info().IID, newDevice(d))
			return oldDevice.Close()
		}

		return nil
	}
	logrus.Info("add device:", d.Info())

	return nil
}

func (m *Manager) CloseDevice(iid string) {
	if v, ok := m.devices.LoadAndDelete(iid); ok {
		err := v.(*device).Close()
		if err != nil {
			logrus.Errorf("manager close device %s err:%v", iid, err)
		}
	}
}

func (m *Manager) IsOTASupport(iid string) (ok bool, err error) {

	d, err := m.GetDeviceInterface(iid)
	if err != nil {
		return
	}

	switch d.(type) {
	case OTADevice:
		return true, nil
	default:
		return
	}
}

func (m *Manager) OTA(iid, firmwareURL string) (ch chan OTAResp, err error) {

	d, err := m.GetDeviceInterface(iid)
	if err != nil {
		return
	}

	switch v := d.(type) {
	case OTADevice:
		return v.OTA(firmwareURL)
	default:
		logrus.Warnf("%s cant't OTA", iid)
		return
	}
}

func (m *Manager) Connect(d *device, params map[string]interface{}) (err error) {

	if ad, authRequired := d.Device.(AuthDevice); authRequired && !ad.IsAuth() {
		if err = ad.Auth(params); err != nil {
			return err
		}
	}

	if v, loaded := m.devices.LoadOrStore(d.Info().IID, d); loaded {
		old := v.(*device)
		old.mutex.Lock()
		if old.connected { // 已连接
			if !old.Online(d.Info().IID) { // 不在线则关闭旧连接
				old.Close()
			} else { // 已经在线则返回
				old.mutex.Unlock()
				return nil
			}
		} else { // 未连接则用旧的未连接设备进行连接
			d = old
		}
		old.mutex.Unlock()
	}
	d.mutex.Lock()
	defer func() {
		d.mutex.Unlock()
		m.discoverDevices.Delete(d.Info().IID)
		// 如果出错，将设备从devices删除
		if err != nil {
			m.devices.Delete(d.Info().IID)
		}
	}()

	if d.connected {
		return
	}
	logrus.Infof("device %s connecting...", d.Info().IID)
	if err = d.Connect(); err != nil {
		return
	}

	logrus.Debugf("device %s connected, define device's thing model", d.Info().IID)
	d.df = definer.NewThingModelDefiner(d.Info().IID, m.notifyAttr, m.notifyThingModelChange)
	err = d.Define(d.df)
	if err != nil {
		if err2 := d.Close(); err2 != nil {
			logrus.Errorf("device close err:%v", err2)
		}
		return
	}
	if d.df != nil {
		d.df.SetNotifyFunc()
		tm := d.df.ThingModel()

		// 记录所有设备（包括子设备）对应的 device
		for _, ins := range tm.Instances {
			m.devices.Store(ins.IID, d)
		}
	}
	d.connected = true

	logrus.Debugf("device %s connect and define done", d.Info().IID)
	return
}

func (m *Manager) Disconnect(iid string, params map[string]interface{}) (err error) {

	defer m.devices.Delete(iid)
	d, err := m.GetDeviceInterface(iid)
	if err != nil {
		return
	}
	if ad, authRequired := d.(AuthDevice); authRequired && ad.IsAuth() {
		if ad.Info().IID == iid {
			if err = ad.RemoveAuthorization(params); err != nil {
				return
			}
		}
	}
	if err = d.Disconnect(iid); err != nil {
		return
	}
	// 子设备不会走发现，删除子设备不要把父设备的连接关掉
	if d.Info().IID == iid {
		err = d.Close()
	}
	return
}

func (m *Manager) HealthCheck(d *device, iid string) bool {
	if d == nil {
		return false
	}

	online := d.Device.Online(iid)
	logrus.Debugf("%s HealthCheck,online: %v", iid, online)
	return online
}

func (m *Manager) notifyEvent(event Event) error {
	if m.eventChannel == nil {
		logrus.Warn("eventChannel not set")
		return nil
	}
	select {
	case m.eventChannel <- event:
	default:
	}

	logrus.Debugf("notifyEvent: %s:%s", event.Type, string(event.Data))
	return nil

}
func (m *Manager) notifyAttr(attrEvent definer.AttributeEvent) (err error) {
	data, _ := json.Marshal(attrEvent)
	ev := Event{
		Type: AttrChangeEvent,
		Data: data,
	}

	return m.notifyEvent(ev)
}

// notifyThingModelAuth d为发现设备
func (m *Manager) notifyThingModelAuth(d *device) (err error) {
	tme := definer.ThingModelEvent{IID: d.Info().IID}
	var ad AuthDevice
	ad, tme.ThingModel.AuthRequired = d.Device.(AuthDevice)
	if !tme.ThingModel.AuthRequired {
		return
	}
	tme.ThingModel.IsAuth = ad.IsAuth()
	tme.ThingModel.AuthParams = ad.AuthParams()
	data, _ := json.Marshal(tme)
	ev := Event{
		Type: ThingModelChangeEvent,
		Data: data,
	}
	return m.notifyEvent(ev)
}

// notifyThingModelChange iid是桥接设备iid，tme.IID是实际更新的设备iid
func (m *Manager) notifyThingModelChange(iid string) (err error) {

	d, err := m.GetDevice(iid)
	if err != nil {
		return
	}
	tme := definer.ThingModelEvent{IID: iid}
	if d.df != nil {
		tme.ThingModel = d.df.ThingModel()
	}
	tme.ThingModel.OTASupport, err = m.IsOTASupport(iid)
	if err != nil {
		return
	}
	var ad AuthDevice
	ad, tme.ThingModel.AuthRequired = d.Device.(AuthDevice)
	if tme.ThingModel.AuthRequired {
		tme.ThingModel.IsAuth = ad.IsAuth()
		tme.ThingModel.AuthParams = ad.AuthParams()
	}
	data, _ := json.Marshal(tme)
	ev := Event{
		Type: ThingModelChangeEvent,
		Data: data,
	}
	// 物模型变更需要重新给新增加的instance设置通知函数
	if d.df != nil {
		d.df.SetNotifyFunc()
	}
	for _, ins := range tme.ThingModel.Instances {
		m.devices.Store(ins.IID, d)
	}
	return m.notifyEvent(ev)
}

func (m *Manager) SetAttributes(as []SetAttribute) error {

	for _, a := range as {
		df, err := m.getDefiner(a.IID)
		if err != nil {
			return err
		}
		if err = df.SetAttribute(a.IID, a.AID, a.Val); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) GetThingModel(iid string) (tm thingmodel.ThingModel, err error) {
	df, err := m.getDefiner(iid)
	if err != nil {
		return
	}
	return df.ThingModel(), nil
}

func (m *Manager) getDefiner(iid string) (df *definer.Definer, err error) {

	d, err := m.GetDevice(iid)
	if err != nil {
		return
	}
	if d.df == nil {
		err = fmt.Errorf("%s definer is nil", iid)
		return
	}
	return d.df, nil
}

func (m *Manager) GetDeviceInterface(iid string) (Device, error) {

	d, err := m.GetDevice(iid)
	if err != nil {
		return nil, err
	}
	return d.Device, nil
}

func (m *Manager) GetDevice(iid string) (*device, error) {

	v, ok := m.devices.Load(iid)
	if !ok {
		return nil, DeviceNotExist
	}
	if d, ok := v.(*device); ok {
		return d, nil
	}
	return nil, fmt.Errorf("%s: is not *device", iid)
}

func (m *Manager) GetDevices() (devices []Device) {
	m.devices.Range(func(key, value interface{}) bool {
		if d, ok := value.(*device); ok {
			devices = append(devices, d.Device)
		}
		return true
	})
	return
}
