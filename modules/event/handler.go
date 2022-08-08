package event

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/zhiting-tech/smartassistant/modules/device"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/plugin"
	"github.com/zhiting-tech/smartassistant/modules/task"
	"github.com/zhiting-tech/smartassistant/modules/websocket"
	_ "github.com/zhiting-tech/smartassistant/modules/websocket"
	"github.com/zhiting-tech/smartassistant/pkg/event"
	"github.com/zhiting-tech/smartassistant/pkg/logger"
	"github.com/zhiting-tech/smartassistant/pkg/plugin/sdk/v2/definer"
	"github.com/zhiting-tech/smartassistant/pkg/thingmodel"
)

func RegisterEventFunc(ws *websocket.Server) {
	event.RegisterEvent(event.AttributeChange, ws.MulticastMsg,
		UpdateDeviceShadowBeforeExecuteTask, RecordDeviceState)
	event.RegisterEvent(event.DeviceDecrease, ws.MulticastMsg)
	event.RegisterEvent(event.DeviceIncrease, ws.MulticastMsg)
	event.RegisterEvent(event.OnlineStatus, ws.MulticastMsg)
	event.RegisterEvent(event.ThingModelChange, UpdateThingModel, ws.MulticastMsg)
}

func UpdateThingModel(em event.EventMessage) (err error) {
	tm := em.Param["thing_model"].(thingmodel.ThingModel)
	areaID := em.AreaID
	pluginID := em.Param["plugin_id"].(string)
	iid := em.Param["iid"].(string)

	// 网关未添加则设备不更新
	var primaryDevice entity.Device
	if primaryDevice, err = entity.GetPluginDevice(areaID, pluginID, iid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	// 更新设备物模型
	if err = primaryDevice.UpdateThingModel(tm); err != nil {
		return
	}

	// 更新子设备物模型
	for _, ins := range tm.GetSubInstances() {

		var e entity.Device
		if e, err = entity.GetPluginDevice(areaID, pluginID, ins.IID); err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}
		// 子设备不存在则为所有角色增加改设备的权限 && 更新子设备房间为网关默认房间
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e, err = plugin.InstanceToEntity(ins, pluginID, iid, areaID)
			if err != nil {
				logrus.Error(err)
				err = nil
				continue
			}

			if err = device.Create(areaID, &e); err != nil {
				logrus.Error(err)
				err = nil
				continue
			}

			// 更新设备房间为默认房间
			updates := map[string]interface{}{
				"location_id":      primaryDevice.LocationID,
				"location_order":   0,
				"department_id":    primaryDevice.DepartmentID,
				"department_order": 0,
			}
			if err = entity.UpdateDeviceWithMap(e.ID, updates); err != nil {
				logrus.Error(err)
				err = nil
				continue
			}
		} else {
			childThingModel := thingmodel.ThingModel{
				Instances:  []thingmodel.Instance{ins},
				OTASupport: false, // TODO 根据插件实现判断
			}
			if err = e.UpdateThingModel(childThingModel); err != nil {
				return
			}
		}

		// 发送通知有设备增加 FIXME 更新也发通知（子设备重复添加时需要通知来完成添加流程）
		m := event.NewEventMessage(event.DeviceIncrease, areaID)
		m.Param = map[string]interface{}{
			"device": e,
		}
		event.Notify(m)

	}

	// 如果有子设备添加到了房间或部门，则重新更新对应房间/部门的设备排序
	if len(tm.GetSubInstances()) != 0 {
		if primaryDevice.LocationID != 0 {
			if err = entity.ReorderAllLocationDevices(primaryDevice.LocationID); err != nil {
				return
			}
		}
		if primaryDevice.DepartmentID != 0 {
			if err = entity.ReorderAllDepartmentDevices(primaryDevice.DepartmentID); err != nil {
				return
			}
		}

	}
	return nil
}

func UpdateDeviceShadowBeforeExecuteTask(em event.EventMessage) (err error) {
	if err = UpdateDeviceShadow(em); err != nil {
		return
	}

	err = ExecuteTask(em)
	return
}

var updateShadowMap sync.Map

func UpdateDeviceShadow(em event.EventMessage) error {

	attr := em.GetAttr()
	if attr == nil {
		logger.Warn(" attr is nil")
		return nil
	}
	deviceID := em.GetDeviceID()

	val, _ := updateShadowMap.LoadOrStore(deviceID, &sync.Mutex{})
	mu := val.(*sync.Mutex)
	mu.Lock()
	defer func() {
		mu.Unlock()
		updateShadowMap.Delete(deviceID)
	}()

	d, err := entity.GetDeviceByID(deviceID)
	if err != nil {
		return err
	}
	// 从设备影子中获取属性
	shadow, err := d.GetShadow()
	if err != nil {
		return err
	}
	shadow.UpdateReported(attr.IID, attr.AID, attr.Val)
	d.Shadow, err = json.Marshal(shadow)
	if err != nil {
		return err
	}
	if err = entity.GetDB().Model(&entity.Device{ID: d.ID}).Update("shadow", d.Shadow).Error; err != nil {
		return err
	}

	return nil
}

func ExecuteTask(em event.EventMessage) error {
	deviceID := em.GetDeviceID()
	d, err := entity.GetDeviceByID(deviceID)
	if err != nil {
		return err
	}
	attr := em.GetAttr()
	if attr == nil {
		logger.Warn("device or attr is nil")
		return nil
	}
	return task.GetManager().DeviceStateChange(d, *attr)
}

type State struct {
	thingmodel.Attribute
}

func RecordDeviceState(em event.EventMessage) (err error) {
	deviceID := em.GetDeviceID()
	d, err := entity.GetDeviceByID(deviceID)
	if err != nil {
		return
	}

	// 将变更的属性保存记录设备状态

	tm, err := d.GetThingModel()
	if err != nil {
		return
	}

	ae := em.Param["attr"].(definer.AttributeEvent)
	attribute, err := tm.GetAttribute(ae.IID, ae.AID)
	if err != nil {
		return
	}
	state := State{Attribute: attribute}
	state.Val = ae.Val
	stateBytes, _ := json.Marshal(state)
	return entity.InsertDeviceState(d, stateBytes)
}
