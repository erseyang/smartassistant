package message

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/device"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
)

type NotificationSettingInfoResp struct {
	Permission bool `json:"permission"`
	Devices    []DeviceSetting `json:"devices"`
}

type DeviceSetting struct {
	Id         int    `json:"id"`
	Permission bool   `json:"permission"`
	Name       string `json:"name"`
	Img        string `json:"img"`
}

func NotificationSettingInfo(c *gin.Context) {
	var (
		err         error
		resp        NotificationSettingInfoResp
		setting     = entity.GetDefaultDevNotificationMessageSetting()
		sessionUser = session.Get(c)
		devices     []entity.Device
		msgSettings []entity.MessageSetting
	)
	defer func() {
		response.HandleResponse(c, err, resp)
	}()

	if err = entity.GetUserSetting(entity.DevNotificationSetting, &setting, sessionUser.AreaID); err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}
	resp.Permission = setting.Permission

	if devices, err = entity.GetDevices(sessionUser.AreaID); err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}

	if msgSettings, err = entity.GetMessageSettings(sessionUser.AreaID, sessionUser.UserID); err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}

	permission := true
	for _, dev := range devices {
		for _, msg := range msgSettings {
			if dev.ID == msg.DeviceId {
				permission = msg.Permission
				break
			}
		}
		resp.Devices = append(resp.Devices, DeviceSetting{
			Id:         dev.ID,
			Permission: permission,
			Name:       dev.Name,
			Img:        device.LogoURL(c.Request, dev),
		})
		permission = true
	}

}
