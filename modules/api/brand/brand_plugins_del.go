package brand

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"
	"github.com/zhiting-tech/smartassistant/pkg/analytics"
	"github.com/zhiting-tech/smartassistant/pkg/event"
)

func DelPlugins(c *gin.Context) {
	var (
		req  handlePluginsReq
		resp handlePluginsResp
		err  error
	)

	defer func() {
		response.HandleResponse(c, err, resp)
	}()
	if err = c.BindUri(&req); err != nil {
		return
	}
	if err = c.BindJSON(&req); err != nil {
		return
	}
	plgs, err := req.GetPluginsWithContext(c.Request.Context())
	if err != nil {
		return
	}

	user := session.Get(c)
	var devices []entity.Device
	for _, plg := range plgs {
		if plg.Brand != req.BrandName {
			continue
		}

		dList, err := entity.GetDevicesByPluginID(plg.ID)
		if err != nil {
			return
		}

		if err = plg.Remove(c.Request.Context()); err != nil {
			return
		}

		// 已删除的设备iid和pluginID
		devices = append(devices, dList...)

		go analytics.RecordStruct(analytics.EventTypePluginDelete, user.UserID, plg)
		resp.SuccessPlugins = append(resp.SuccessPlugins, plg.ID)
	}

	em := event.NewEventMessage(event.DeviceDecrease, session.Get(c).AreaID)
	em.Param["deleted_devices"] = devices
	event.Notify(em)
	return
}
