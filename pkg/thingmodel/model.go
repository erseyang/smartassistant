package thingmodel

import (
	"errors"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type ThingModel struct {
	Instances []Instance `json:"instances"`

	// OTASupport 是否支持OTA
	OTASupport bool `json:"ota_support"`

	// AuthRequired 是否需要认证
	AuthRequired bool `json:"auth_required"`
	// IsAuth 是否已经认证
	IsAuth bool `json:"is_auth"`
	// AuthParams 认证需要的参数
	AuthParams []AuthParam `json:"auth_params"`
}

type AuthParam struct {
	Name     string `json:"name"`
	Type     string `json:"type"` // string/bool/int/float/select
	Required bool   `json:"required"`

	Default interface{} `json:"default,omitempty"`
	Min     interface{} `json:"min,omitempty"`
	Max     interface{} `json:"max,omitempty"`
	Options []Option    `json:"options,omitempty"`
}
type Option struct {
	Name string      `json:"name"`
	Val  interface{} `json:"val"`
}

func (das ThingModel) GetAttribute(iid string, aid int) (Attribute, error) {

	ins, err := das.GetInstance(iid)
	if err != nil {
		return Attribute{}, err
	}
	return ins.GetAttribute(aid)
}

func (das ThingModel) GetInstance(iid string) (Instance, error) {

	for _, ins := range das.Instances {
		if ins.IID == iid {
			return ins, nil
		}
	}
	return Instance{}, fmt.Errorf("instance %s not found", iid)
}

// GetInfo 获取设备基础信息
func (das ThingModel) GetInfo(iid string) (info DeviceInfo, err error) {

	ins, err := das.GetInstance(iid)
	if err != nil {
		return
	}
	return ins.GetInfo()
}

// IsBridge 是否是桥接类设备，如网关、无线控制器
func (das ThingModel) IsBridge() bool {
	for _, i := range das.Instances {
		if i.IsBridge() {
			return true
		}
	}
	return false
}

// BridgeInstance 获取桥接类设备的instance
func (das ThingModel) BridgeInstance() (Instance, error) {
	for _, i := range das.Instances {
		if i.IsBridge() {
			return i, nil
		}
	}
	return Instance{}, fmt.Errorf("device is not a bridge")
}

// PrimaryInstance 获取设备的主要instance，就是设备本身的instance
func (das ThingModel) PrimaryInstance() (Instance, error) {
	if len(das.Instances) == 0 {
		return Instance{}, fmt.Errorf("no instance in thing model")
	}
	if len(das.Instances) == 1 {
		return das.Instances[0], nil
	}
	// 网关等多instance的设备返回交接设备的instance
	return das.BridgeInstance()
}

func (das ThingModel) GetSubInstances() (subInstance []Instance) {
	if len(das.Instances) == 1 {
		return
	}
	for _, i := range das.Instances {
		if i.IsBridge() {
			continue
		}
		subInstance = append(subInstance, i)
	}
	return
}

type DeviceInfo struct {
	IID          string
	Name         string
	Model        string
	Manufacturer string
	Version      string
	Type         string
}

type Instance struct {
	IID      string    `json:"iid"`
	Services []Service `json:"services"`

	sync.Mutex
}

func (i Instance) GetInfo() (info DeviceInfo, err error) {

	for _, srv := range i.Services {
		if srv.Type != "info" {
			continue
		}
		for _, attr := range srv.Attributes {
			switch attr.Type {
			case Name.String():
				info.Name, _ = attr.Val.(string)
			case Model.String():
				info.Model, _ = attr.Val.(string)
			case Manufacturer.String():
				info.Manufacturer, _ = attr.Val.(string)
			case Identify.String():
				info.IID, _ = attr.Val.(string)
			case Version.String():
				info.Version, _ = attr.Val.(string)
			case Type.String():
				info.Type, _ = attr.Val.(string)
			}
		}

		return
	}
	err = errors.New("info service not found")
	return
}

// IsBridge 是否是桥接类设备
func (i Instance) IsBridge() bool {
	for _, s := range i.Services {
		if s.Type == GatewayService {
			return true
		}
	}
	return false
}

func (i Instance) GetAttribute(aid int) (Attribute, error) {
	for _, s := range i.Services {
		for _, a := range s.Attributes {
			if a.AID == aid {
				return a, nil
			}
		}
	}
	return Attribute{}, fmt.Errorf("attribute %d not found", aid)
}

func (i *Instance) UpdateServiceName(index int, name string) error {
	if index >= len(i.Services) {
		return fmt.Errorf("%s services[%d] not found", i.IID, index)
	}
	i.Services[index].Name = name
	return nil
}

func (tm *ThingModel) Update(update ThingModel) (err error) {
	for _, updateInstance := range update.Instances {
		var found bool
		for i, _ := range tm.Instances {
			if tm.Instances[i].IID == updateInstance.IID {
				found = true
				if err = tm.Instances[i].Update(updateInstance); err != nil {
					return
				}
			}
		}
		if !found {
			tm.Instances = append(tm.Instances, updateInstance)
		}
	}
	tm.OTASupport = update.OTASupport
	tm.AuthRequired = update.AuthRequired
	tm.AuthParams = update.AuthParams
	tm.IsAuth = update.IsAuth
	return nil
}

func (tm *ThingModel) UpdateServiceName(iid string, index int, name string) error {
	for _, ins := range tm.Instances {
		if ins.IID == iid {
		}
		return ins.UpdateServiceName(index, name)
	}
	return nil
}

// Update service只处理增加的情况
func (old *Instance) Update(update Instance) (err error) {

	// 防止并发更新services
	old.Lock()
	defer old.Unlock()
	length := len(old.Services)
	for j, updateService := range update.Services {
		if j < length {
			if err = old.Services[j].Update(updateService); err != nil {
				return
			}
		} else {
			old.Services = append(old.Services, updateService)
		}
	}
	return
}

// Update 仅更新所有属性，不会更新默认名字和类型
func (old *Service) Update(update Service) (err error) {
	if old.Type != update.Type {
		// 忽略非同类型的service的更新
		logrus.Warnf("service type mismatch: %s!=%s", old.Type, update.Type)
		return nil
	}
	old.Attributes = make([]Attribute, len(update.Attributes))
	copy(old.Attributes, update.Attributes)
	return nil
}
