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

type SilenceDevReq struct {
	DeviceId int  `json:"id" uri:"id"`
	State    bool `json:"state"`
}

func UpdateSilenceDev(c *gin.Context) {
	var (
		err error
		req SilenceDevReq

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
		return
	}

	setting := entity.MessageSetting{
		UserId:   sessionUser.UserID,
		AreaId:   sessionUser.AreaID,
		DeviceId: req.DeviceId,
		State:    req.State,
	}

	if err = entity.UpdateMessageSettingState(setting); err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}
}

func (req SilenceDevReq) validate() (err error) {
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
