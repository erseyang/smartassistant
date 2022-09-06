package message

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/device"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
)

type SilenceDevResp struct {
	Devices []Device `json:"devices"`
}

type Device struct {
	Id    int    `json:"id"`
	State bool   `json:"state"`
	Name  string `json:"name"`
	Img   string `json:"img"`
}

func SilenceDevList(c *gin.Context) {
	var (
		err         error
		resp        SilenceDevResp
		devices     []entity.Device
		msgSettings []entity.MessageSetting
		sessionUser = session.Get(c)
	)
	defer func() {
		response.HandleResponse(c, err, resp)
	}()

	if devices, err = entity.GetDevices(sessionUser.AreaID); err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}

	if msgSettings, err = entity.GetMessageSettings(sessionUser.AreaID, sessionUser.UserID); err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}

	state := true
	for _, dev := range devices {
		for _, msg := range msgSettings {
			if dev.ID == msg.DeviceId {
				state = msg.State
				break
			}
		}
		resp.Devices = append(resp.Devices, Device{
			Id:    dev.ID,
			State: state,
			Name:  dev.Name,
			Img:   device.LogoURL(c.Request, dev),
		})
		state = true
	}
}
