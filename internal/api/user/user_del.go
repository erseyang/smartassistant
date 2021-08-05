package user

import (
	"github.com/zhiting-tech/smartassistant/internal/api/utils/cloud"
	"github.com/zhiting-tech/smartassistant/internal/api/utils/response"
	"github.com/zhiting-tech/smartassistant/internal/entity"
	"github.com/zhiting-tech/smartassistant/internal/types/status"
	"github.com/zhiting-tech/smartassistant/internal/utils/session"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
)

// DelUser 用于处理删除成员接口的请求
func DelUser(c *gin.Context) {
	var (
		err         error
		userID      int
		sessionUser *session.User
	)

	defer func() {
		response.HandleResponse(c, err, nil)
	}()

	userID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		err = errors.Wrap(err, errors.BadRequest)
		return
	}

	sessionUser = session.Get(c)
	if sessionUser == nil {
		err = errors.Wrap(err, status.AccountNotExistErr)
		return
	}

	if _, err = entity.GetUserByID(userID); err != nil {
		return
	}

	// 成员本人不能删除自己
	if sessionUser.UserID == userID {
		err = errors.Wrap(err, status.DelSelfErr)
		return
	}

	if entity.IsSAOwner(userID) {
		err = errors.New(status.Deny)
		return
	}

	if err = entity.DelUser(userID); err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
	}
	cloud.RemoveSAUser(userID)
	return

}
