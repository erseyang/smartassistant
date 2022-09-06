package message

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/device"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/types/status"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"
	"github.com/zhiting-tech/smartassistant/modules/utils/url"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
	"github.com/zhiting-tech/smartassistant/pkg/logger"
)

type MessagesResp struct {
	Messages []NotificationMessage `json:"messages"`
}

type MessagesReq struct {
	Type  int `json:"type" form:"type"`
	MsgId int `json:"msg_id" form:"msg_id"`
}

func GetMessages(c *gin.Context) {
	var (
		resp       MessagesResp
		err        error
		req        MessagesReq
		user       = session.Get(c)
		msgRecords []entity.MessageRecordInfo
	)

	defer func() {
		response.HandleResponse(c, err, resp)
	}()

	if err = c.Bind(&req); err != nil {
		err = errors.Wrap(err, errors.BadRequest)
		return
	}

	if err = req.validate(); err != nil {
		return
	}

	if req.MsgId != 0 {
		if msgRecords, err = entity.MessageRecordInfoList(req.MsgId, user.UserID, 10, req.Type, user.AreaID); err != nil {
			err = errors.Wrap(err, errors.InternalServerErr)
			return
		}
	} else {
		if msgRecords, err = entity.MessageRecordLatestList(user.UserID, 10, req.Type, user.AreaID); err != nil {
			err = errors.Wrap(err, errors.InternalServerErr)
			return
		}
	}

	setting := entity.GetDefaultMessageReadMsgId()
	err = entity.GetUserSetting(entity.MsgCenterReadMsgId, &setting, user.AreaID)
	if err != nil {
		logger.Warning("GetUserSetting ReadMsgId error:", err)
	}

	if len(msgRecords) != 0 {
		if msgRecords[0].ID > setting.MsgId{
			setting.MsgId = msgRecords[0].ID
			err = entity.UpdateUserSetting(entity.MsgCenterReadMsgId, &setting, user.UserID, user.AreaID)
			if err != nil {
				logger.Warning("UpdateUserSetting ReadMsgId error:", err)
			}
		}
	}

	for _, msg := range msgRecords {

		notifyMsg := NotificationMessage{
			ID:          msg.ID,
			Type:        msg.Type,
			Location:    msg.LocationName,
			CreatedTime: msg.CreatedAt.Unix(),
			Title:       msg.Title,
			Content:     msg.Content,
		}
		if notifyMsg.Type == entity.NotifyTypeAlarm {
			notifyMsg.LogoUrl = device.LogoURL(c.Request, entity.Device{Type: msg.DeviceType, Model: msg.Model, PluginID: msg.PluginId})
		} else if notifyMsg.Type == entity.NotifyTypeArea {
			notifyMsg.LogoUrl = url.ImageUrl(c.Request, entity.AreaLogoUrlMap[notifyMsg.Title])
		}
		resp.Messages = append(resp.Messages, notifyMsg)
	}
}

func (req MessagesReq) validate() (err error) {
	if req.Type != entity.NotifyTypeAlarm && req.Type != entity.NotifyTypeArea && req.Type != 0 {
		err = errors.Wrap(err, status.MessageCenterTypeNoExist)
	}
	return
}
