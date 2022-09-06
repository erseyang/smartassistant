package message

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/types/status"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
)

type NotificationReq struct {
	Type    int    `json:"type"` // 类型，1 第三方授权成功
	UserId  int    `json:"user_id" form:"user_id"`
	AreaId  uint64 `json:"area_id" form:"area_id"`
	AppName string `json:"app_name" form:"app_name"`
}

func Notification(c *gin.Context) {
	var (
		req       NotificationReq
		err       error
		user      entity.User
		msgRecord entity.MessageRecord
	)
	defer func() {
		response.HandleResponse(c, err, nil)
	}()

	if err = c.BindJSON(&req); err != nil {
		err = errors.Wrap(err, errors.BadRequest)
		return
	}

	if err = req.validate(); err != nil {
		return
	}

	if user, err = entity.GetUserByID(req.UserId); err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}

	switch req.Type {
	case TypePlatformAuth:
		msgRecord = entity.NewAreaMessageRecord(req.AreaId, entity.TitlePlatformChange, fmt.Sprintf(ContentPlatformAuth, user.Nickname, req.AppName))
	}

	// 入库推送消息
	go GetMessagesManager().SendMsg(msgRecord)

}

func (req NotificationReq) validate() (err error) {
	if req.Type != TypePlatformAuth {
		err = errors.Wrap(err, status.MessageCenterTypeNotExistErr)
		return
	}
	return
}
