package area

import (
	"fmt"
	"github.com/zhiting-tech/smartassistant/modules/api/message"
	"github.com/zhiting-tech/smartassistant/modules/extension"
	pb "github.com/zhiting-tech/smartassistant/pkg/extension/proto"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/cloud"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/types/status"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
)

// QuitArea 用于处理退出家庭接口的请求
func QuitArea(c *gin.Context) {
	var (
		err         error
		sessionUser *session.User
		userID      int
		areaID      uint64
		area        entity.Area
	)

	defer func() {
		response.HandleResponse(c, err, nil)
	}()

	sessionUser = session.Get(c)
	if sessionUser == nil {
		err = errors.Wrap(err, status.AccountNotExistErr)
		return
	}

	if areaID, err = strconv.ParseUint(c.Param("id"), 10, 64); err != nil {
		err = errors.Wrap(err, errors.BadRequest)
		return
	}
	if area, err = entity.GetAreaByID(areaID); err != nil {
		return
	}

	if entity.IsOwnerOfArea(sessionUser.UserID, areaID) {
		areaTypeStr := area.AreaType.String()
		err = errors.Wrapf(err, status.OwnerQuitErr, areaTypeStr, areaTypeStr)
		return
	}

	if userID, err = strconv.Atoi(c.Param("user_id")); err != nil {
		err = errors.Wrap(err, errors.BadRequest)
		return
	}

	if userID != sessionUser.UserID {
		err = errors.New(status.Deny)
		return
	}

	// 退出家庭删除网盘所有文件夹
	extension.GetExtensionServer().Notify(pb.SAEvent_del_user_ev, map[string]interface{}{
		"ids": []int{userID},
	})

	if err = cloud.RemoveSAUserWithContext(c.Request.Context(), areaID, userID); err != nil {
		return
	}

	if err = entity.DelUser(sessionUser.UserID); err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}

	user, _ := entity.GetUserByID(userID)
	// 入库推送成员退出消息
	msgRecord := entity.NewMemberChangeMessageRecord(areaID, fmt.Sprintf(message.ContentMemberExit, user.Nickname, area.AreaType.String()))
	go message.GetMessagesManager().SendMsg(msgRecord)

	return

}
