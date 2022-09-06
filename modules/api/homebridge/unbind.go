package homebridge

import (
	errors2 "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/zhiting-tech/smartassistant/modules/api/message"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/types/status"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
	"github.com/zhiting-tech/smartassistant/pkg/logger"
	"net/http"
)

func Unbind(c *gin.Context) {
	var err error
	defer func() {
		response.HandleResponse(c, err, nil)
	}()

	if err = UnbindHomeBridge(c); err != nil {
		return
	}

	sessionUser := session.Get(c)
	user, err := entity.GetUserByID(sessionUser.UserID)
	if err != nil {
		logger.Warning("homebridge SaveAndSendMsg err", err)
		return
	}

	// 入库推送家居桥接解除授权消息
	msgRecord := entity.NewHomeBridgeChangeMessageRecord(sessionUser.AreaID, fmt.Sprintf(message.ContentHomeBridgeUnbind, user.Nickname))
	go message.GetMessagesManager().SendMsg(msgRecord)
}

func UnbindHomeBridge(c *gin.Context) (err error) {
	b, err := reqToHomeBridge(c, http.MethodDelete, getHomeBridgeApi("unbind"), nil)
	if err != nil {
		err = errors.Wrap(err, status.UnbindError)
		return
	}

	code := gjson.GetBytes(b, "status").Int()
	if code != 0 {
		e := errors2.New(gjson.GetBytes(b, "reason").String())
		err = errors.Wrap(e, status.UnbindError)
		return
	}
	return
}
