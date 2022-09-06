package message

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
	"gorm.io/gorm"
)

type GetReadStatusResp struct {
	IsRead int `json:"is_read"`
}

func GetReadStatus(c *gin.Context) {
	var (
		err         error
		resp        GetReadStatusResp
		sessionUser = session.Get(c)
		msgRecord   entity.MessageRecord
	)
	defer func() {
		response.HandleResponse(c, err, resp)
	}()

	if msgRecord, err = entity.GetLatestMsgRecord(sessionUser.UserID, sessionUser.AreaID); err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.IsRead = 1
			err = nil
			return
		}
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}

	setting := entity.GetDefaultMessageReadMsgId()
	err = entity.GetUserSetting(entity.MsgCenterReadMsgId, &setting, sessionUser.AreaID)
	if err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}
	if msgRecord.ID > setting.MsgId{
		resp.IsRead = 0
	}else {
		resp.IsRead = 1
	}
}
