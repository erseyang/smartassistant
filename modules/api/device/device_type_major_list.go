package device

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/plugin"
)

type ModelDevice struct {
	Name         string `json:"name"`
	Model        string `json:"model"`
	Manufacturer string `json:"manufacturer"`
	Logo         string `json:"logo"`         // logo地址
	Provisioning string `json:"provisioning"` // 配置页地址
	PluginID     string `json:"plugin_id"`

	Protocol string `json:"protocol"` // 连接云端的协议类型，tcp/mqtt
}

type Type struct {
	Name string            `json:"name"`
	Type plugin.DeviceType `json:"type"`
}

type Types []Type

type MajorResp struct {
	Types `json:"types"`
}

var majorTypes = []Type{
	{"照明", plugin.TypeLight},
	{"插座", plugin.TypeOutlet},
	{"开关", plugin.TypeSwitch},
	{"传感器", plugin.TypeSensor},
	{"路由网关", plugin.TypeRoutingGateway},
	{"安防", plugin.TypeSecurity},
	{"生活电器", plugin.TypeLifeElectric},
}

func getParentType(parentType plugin.DeviceType) (t Type, ok bool) {
	for _, m := range majorTypes {
		if m.Type == parentType {
			return m, true
		}
	}
	return Type{}, false
}

// MajorTypeList 获取主分类
func MajorTypeList(c *gin.Context) {
	var (
		err  error
		resp MajorResp
	)
	resp.Types = make([]Type, 0)

	defer func() {
		response.HandleResponse(c, err, resp)
	}()

	configs := plugin.GetGlobalClient().Configs()
	m := make(map[plugin.DeviceType]Type, 0)
	for _, plgConf := range configs {
		for _, d := range plgConf.SupportDevices {
			if d.Provisioning == "" { // 没有配置置网页则忽略
				continue
			}

			pType, ok := getCurrentType(d.Type)
			if !ok || pType.ParentType == "" {
				continue
			}
			if v, ok := getParentType(pType.ParentType); ok {
				m[pType.ParentType] = v
			}
		}
	}

	for _, mt := range majorTypes {
		if _, ok := m[mt.Type]; ok {
			resp.Types = append(resp.Types, mt)
		}
	}
}
