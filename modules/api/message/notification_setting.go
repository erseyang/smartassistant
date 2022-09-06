package message

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
)

type NotificationSettingReq struct {
	Permission bool `json:"permission"`
}

func UpdateNotificationSetting(c *gin.Context) {
	var (
		req         NotificationSettingReq
		err         error
		sessionUser *session.User
	)

	defer func() {
		response.HandleResponse(c, err, nil)
	}()

	if err = c.BindJSON(&req); err != nil {
		err = errors.Wrap(err, errors.BadRequest)
		return
	}

	// 只有SA拥有者才能设置
	sessionUser = session.Get(c)

	err = req.UpdateNotificationSetting(sessionUser.UserID, sessionUser.AreaID)
	if err != nil {
		return
	}

}
func (req *NotificationSettingReq) UpdateNotificationSetting(userId int, areaID uint64)(err error) {
	setting := entity.MessageDevNotificationSetting{
		Permission: req.Permission,
	}
	err = entity.UpdateUserSetting(entity.DevNotificationSetting, &setting, userId, areaID)

	if err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}

	return nil
}
