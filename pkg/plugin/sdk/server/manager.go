package server

import (
	"errors"
	"log"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/zhiting-tech/smartassistant/pkg/plugin/sdk/attribute"
	"github.com/zhiting-tech/smartassistant/pkg/plugin/sdk/utils"
)

type Manager struct {
	devices    sync.Map
	notifyChan chan Notify

	notifyChans map[chan Notify]struct{}
}

func NewManager() *Manager {
	return &Manager{}
}

func (p *Manager) init() {
	p.notifyChans = make(map[chan Notify]struct{})
	p.notifyChan = make(chan Notify, 10)
	go func() {
		for {
			select {
			case n := <-p.notifyChan:
				for ch := range p.notifyChans {
					select {
					case ch <- n:
					default:
					}
				}
			}
		}
	}()
}

func (p *Manager) setAttributeNotify(identity string) error {
	device, ok := p.devices.Load(identity)
	if !ok {
		return errors.New("setAttributeNotify error,no device found")
	}

	s := utils.Parse(device)
	for _, instance := range s.Instances {
		for _, attr := range instance.Attributes {
			// TODO 忽略未设置属性
			if n, ok := attr.Model.(attribute.Notifier); ok {
				n.SetNotifyFunc(p.Notify(identity, instance.ID, attr))
			}
		}
	}
	return nil
}
func (p *Manager) AddDevice(device Device) error {
	if device == nil {
		return errors.New("device is nil")
	}
	_, loaded := p.devices.LoadOrStore(device.Identity(), device)
	if loaded {
		logrus.Debug("device already exist")
		device.Close()
		return nil
	}
	if err := device.Setup(); err != nil {
		return err
	}
	logrus.Info("add device:", device.Info())

	go p.WatchNotify(device)
	return p.setAttributeNotify(device.Identity())
}

func (p *Manager) WatchNotify(device Device) {
	s := utils.Parse(device)
	ch := device.GetChannel()

	for {
		select {
		case v, ok := <-ch:
			if !ok {
				err := errors.New("device channel close")
				logrus.Error(err)
				return
			}
			attr := s.GetAttribute(v.InstanceID, v.Attr)
			if attr == nil {
				continue
			}
			if notifier, ok := attr.Model.(attribute.Notifier); ok {
				if err := notifier.Notify(v.Val); err != nil {
					logrus.Error(err)
					return
				}
			}
		}
	}
}

func (p *Manager) Devices() (ds []Device, err error) {
	p.devices.Range(func(key, value interface{}) bool {
		d := value.(Device)
		ds = append(ds, d)
		return true
	})
	return
}

func (p *Manager) Subscribe(notify chan Notify) {
	p.notifyChans[notify] = struct{}{}
}

func (p *Manager) Unsubscribe(notify chan Notify) {
	delete(p.notifyChans, notify)
}

func (p *Manager) Notify(identity string, instanceID int, attr *utils.Attribute) attribute.NotifyFunc {
	return func(val interface{}) error {
		if p.notifyChan == nil {
			logrus.Warn("notifyChan not set")
			return nil
		}
		n := Notify{Identity: identity, InstanceID: instanceID}
		n.Attribute = Attribute{
			ID:        attr.ID,
			Attribute: attr.Name,
			Val:       val,
			ValType:   attr.Type,
		}
		if num, ok := attr.Model.(attribute.IntType); ok {
			n.Attribute.Min, n.Attribute.Max = num.GetRange()
		}
		select {
		case p.notifyChan <- n:
		default:
		}

		log.Println("notify", identity, instanceID, attr, val)
		return nil
	}
}
func (p *Manager) GetAttributes(identity string) (s []Instance, err error) {

	v, ok := p.devices.Load(identity)
	if !ok {
		err = errors.New("device not exist")
		return
	}
	device := v.(Device)
	if err = device.Update(); err != nil { // update value
		return
	}
	return p.getInstances(device), nil
}

func (p *Manager) getInstances(device Device) (instances []Instance) {

	// parse device
	d := utils.Parse(device)
	logrus.Debugf("total %d instances\n", len(d.Instances))
	for _, ins := range d.Instances {

		var attrs []Attribute
		logrus.Debugf("total %d attrs of instance %d\n", len(ins.Attributes), ins.ID)
		for _, attr := range ins.Attributes {
			a := Attribute{
				ID:        attr.ID,
				Attribute: attr.Name,
				Val:       attribute.ValueOf(attr.Model),
				ValType:   attr.Type,
			}
			if num, ok := attr.Model.(attribute.IntType); ok {
				a.Min, a.Max = num.GetRange()
			}

			attrs = append(attrs, a)
		}

		instance := Instance{
			Type:       ins.Type,
			InstanceId: ins.ID,
			Attributes: attrs,
		}
		instances = append(instances, instance)
	}
	return
}
func (p *Manager) SetAttribute(identity string, instanceID int, attr string, val interface{}) (err error) {

	v, ok := p.devices.Load(identity)
	if !ok {
		err = errors.New("device not exist")
		return
	}

	device := v.(Device)
	// parse device
	d := utils.Parse(device)
	a := d.GetAttribute(instanceID, attr)
	if a != nil {
		if setter, ok := a.Model.(attribute.Setter); ok {
			return setter.Set(val)
		}
		return errors.New("attribute not setter")
	}
	return errors.New("instance not found")
}
