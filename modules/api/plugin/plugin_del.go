package plugin

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/plugin"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"
	"github.com/zhiting-tech/smartassistant/pkg/event"
)

type delPluginReq struct {
	PluginID string `uri:"id"`
}

func DelPlugin(c *gin.Context) {
	var (
		err  error
		req  delPluginReq
		resp interface{}
	)

	defer func() {
		response.HandleResponse(c, err, resp)
	}()

	if err = c.BindUri(&req); err != nil {
		return
	}

	sessionUser := session.Get(c)
	p, err := entity.GetPlugin(req.PluginID, sessionUser.AreaID)
	if err != nil {
		return
	}

	dList, err := entity.GetDevicesByPluginID(p.PluginID)
	if err != nil {
		return
	}

	plg := plugin.NewFromEntity(p)
	if err = plg.Remove(c.Request.Context()); err != nil {
		return
	}

	em := event.NewEventMessage(event.DeviceDecrease, sessionUser.AreaID)
	em.Param["deleted_device"] = dList
	event.Notify(em)
}
