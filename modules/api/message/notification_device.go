package message

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/types/status"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
	"gorm.io/gorm"
)

type NotificationDevReq struct {
	DeviceId   int  `json:"id" uri:"id"`
	Permission bool `json:"permission"`
}

func UpdateNotificationDev(c *gin.Context) {
	var (
		err error
		req NotificationDevReq

		sessionUser = session.Get(c)
	)
	defer func() {
		response.HandleResponse(c, err, nil)
	}()

	if err = c.BindUri(&req); err != nil {
		err = errors.Wrap(err, errors.BadRequest)
		return
	}

	if err = c.BindJSON(&req); err != nil {
		err = errors.Wrap(err, errors.BadRequest)
		return
	}

	if err = req.validate(); err != nil {
		err = errors.Wrap(err, errors.BadRequest)
		return
	}

	setting := entity.MessageSetting{
		UserId:     sessionUser.UserID,
		AreaId:     sessionUser.AreaID,
		DeviceId:   req.DeviceId,
		Permission: req.Permission,
	}

	if err = entity.UpdateMessageSettingPermission(setting); err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}
}

func (req NotificationDevReq) validate() (err error) {
	_, err = entity.GetDeviceByID(req.DeviceId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.Wrap(err, status.MessageCenterDeviceNoExist)
		} else {
			err = errors.Wrap(err, errors.InternalServerErr)
		}
	}
	return
}
