package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/message"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/cloud"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/types/status"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
	"github.com/zhiting-tech/smartassistant/pkg/logger"
)

type UnbindAppReq struct {
	AppID  int    `uri:"id"`
	AreaID uint64 `uri:"area_id"`
}

// UnbindApp 用于处理解绑云端接口的请求
func UnbindApp(c *gin.Context) {
	var (
		req     UnbindAppReq
		err     error
		appName string
	)

	defer func() {
		response.HandleResponse(c, err, nil)
	}()

	if err = c.BindUri(&req); err != nil {
		err = errors.Wrap(err, errors.BadRequest)
		return
	}

	u := session.Get(c)
	if err = req.validate(u); err != nil {
		return
	}

	apps, err := cloud.GetAppList(c.Request.Context(), u.AreaID)
	if err != nil {
		return
	}

	if err = cloud.UnbindApp(c.Request.Context(), u.AreaID, req.AppID); err != nil {
		return
	}

	for _, app := range apps {
		if app.AppID == req.AppID {
			appName = app.Name
		}
	}

	user, err := entity.GetUserByID(u.UserID)
	if err != nil {
		logger.Warning("app SaveAndSendMsg err", err)
		return
	}
	// 入库推送第三方平台解除授权消息
	msgRecord := entity.NewPlatformChangeMessageRecord(req.AreaID, fmt.Sprintf(message.ContentPlatformUnbind, user.Nickname, appName))
	go message.GetMessagesManager().SendMsg(msgRecord)
}

func (req UnbindAppReq) validate(u *session.User) (err error) {
	if !u.IsOwner {
		err = errors.New(status.Deny)
		return
	}

	if req.AreaID != u.AreaID {
		err = errors.New(status.Deny)
		return
	}

	return
}
