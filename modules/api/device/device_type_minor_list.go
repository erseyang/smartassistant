package device

import (
	"github.com/gin-gonic/gin"

	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/plugin"
	"github.com/zhiting-tech/smartassistant/modules/types/status"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
	"github.com/zhiting-tech/smartassistant/pkg/logger"
)

type MinorResp struct {
	Types MinorTypes `json:"types"`
}

type MinorTypes []MinorType

type MinorType struct {
	Type
	Devices []ModelDevice `json:"devices"`
}

type MinorReq struct {
	Type plugin.DeviceType `form:"type"`
}

type minorType struct {
	Name       string
	ParentType plugin.DeviceType
	CurType    plugin.DeviceType
}

var minorTypes = []minorType{
	{"台灯", plugin.TypeLight, plugin.TypeLamp},
	{"吸顶灯", plugin.TypeLight, plugin.TypeCeilingLamp},
	{"灯泡", plugin.TypeLight, plugin.TypeBulb},
	{"灯带", plugin.TypeLight, plugin.TypeLightStrip},
	{"吊灯", plugin.TypeLight, plugin.TypePendantLight},
	{"床头灯", plugin.TypeLight, plugin.TypeNightLight},
	{"夜灯", plugin.TypeLight, plugin.TypeBedSideLamp},
	{"风扇灯", plugin.TypeLight, plugin.TypeFanLamp},
	{"筒射灯", plugin.TypeLight, plugin.TypeDownLight},
	{"磁吸轨道灯", plugin.TypeLight, plugin.TypeMagneticRailLamp},

	{"排插", plugin.TypeOutlet, plugin.TypePowerStrip},
	{"转换器", plugin.TypeOutlet, plugin.TypeConverter},
	{"入墙插座", plugin.TypeOutlet, plugin.TypeWallPlug},

	{"单键开关", plugin.TypeSwitch, plugin.TypeOneKeySwitch},
	{"双键开关", plugin.TypeSwitch, plugin.TypeTwoKeySwitch},
	{"三键开关", plugin.TypeSwitch, plugin.TypeThreeKeySwitch},
	{"四键开关", plugin.TypeSwitch, plugin.TypeFourKeySwitch},
	{"无线开关", plugin.TypeSwitch, plugin.TypeWirelessSwitch},
	{"控制器", plugin.TypeSwitch, plugin.TypeController},

	{"温湿度传感器", plugin.TypeSensor, plugin.TypeTemperatureAndHumiditySensor},
	{"人体传感器", plugin.TypeSensor, plugin.TypeHumanSensors},
	{"烟雾传感器", plugin.TypeSensor, plugin.TypeSmokeSensor},
	{"燃气传感器", plugin.TypeSensor, plugin.TypeGasSensor},
	{"门窗传感器", plugin.TypeSensor, plugin.TypeWindowDoorSensor},
	{"水浸传感器", plugin.TypeSensor, plugin.TypeWaterLeakSensor},
	{"光照度传感器", plugin.TypeSensor, plugin.TypeIlluminanceSensor},
	{"动静传感器", plugin.TypeSensor, plugin.TypeDynamicAndStaticSensor},

	{"路由器", plugin.TypeRoutingGateway, plugin.TypeRouter},
	{"网关", plugin.TypeRoutingGateway, plugin.TypeGateway},
	{"Wi-Fi信号放大器", plugin.TypeRoutingGateway, plugin.TypeWifiRepeater},

	{"摄像头", plugin.TypeSecurity, plugin.TypeCamera},
	{"猫眼门铃", plugin.TypeSecurity, plugin.TypePeepholeDoorbell},
	{"门锁", plugin.TypeSecurity, plugin.TypeDoorLock},

	{"窗帘电机", plugin.TypeLifeElectric, plugin.TypeCurtain},
}

func getCurrentType(currentType plugin.DeviceType) (t minorType, ok bool) {
	for _, m := range minorTypes {
		if m.CurType == currentType {
			return m, true
		}
	}
	return minorType{}, false
}

// MinorTypeList 根据主分类获取次级分类和设备类型
func MinorTypeList(c *gin.Context) {
	var (
		err  error
		resp MinorResp
		req  MinorReq
	)
	resp.Types = make(MinorTypes, 0)

	defer func() {
		response.HandleResponse(c, err, resp)
	}()

	if err = c.BindQuery(&req); err != nil {
		err = errors.Wrap(err, errors.BadRequest)
		return
	}
	if _, ok := getParentType(req.Type); !ok {
		err = errors.Wrap(err, status.DeviceTypeNotExist)
		return
	}

	var token string
	u := session.Get(c)
	if u != nil {
		token = u.Token
	}
	pluginConfigs := plugin.GetGlobalClient().Configs()
	m := make(map[plugin.DeviceType][]ModelDevice)
	for _, pluginConf := range pluginConfigs {

		for _, d := range pluginConf.SupportDevices {
			if d.Provisioning == "" { // 没有配置置网页则忽略
				continue
			}
			// 拼接token和插件id辅助插件实现websocket请求
			provisioning, err := d.WrapProvisioning(token, pluginConf.ID)
			if err != nil {
				logger.Error(err)
				continue
			}
			md := ModelDevice{
				Name:         d.Name,
				Manufacturer: pluginConf.Brand,
				Model:        d.Model,
				Logo: plugin.PluginTargetURL(c.Request, pluginConf.ID,
					d.Model, d.Logo), // 根据配置拼接插件中的图片地址
				Provisioning: provisioning,
				PluginID:     pluginConf.ID,
				Protocol:     d.Protocol,
			}
			pType, ok := getCurrentType(d.Type)
			if !ok{
				continue
			}
			if req.Type == pType.ParentType || req.Type == d.Type {
				m[d.Type] = append(m[d.Type], md)
			}
		}
	}
	for _, mt := range minorTypes {
		if val, ok := m[mt.CurType]; ok {
			name := mt.Name
			if mt.Name == "" {
				name = "其他"
			}
			resp.Types = append(resp.Types, MinorType{Type{name, mt.CurType}, val})
		}
	}
}
