package message

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/types/status"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
	"time"
)

type SilenceSettingReq struct {
	State       bool `json:"state"`         // 免打扰状态，true 开启，false 关闭
	IsAllDevice bool `json:"is_all_device"` // 是否为所有设备，ture 是，false 否
	// 生效时间的配置
	EffectStart int64 `json:"effect_start"`
	EffectEnd   int64 `json:"effect_end"`

	// 重复执行的配置
	RepeatType entity.RepeatType `json:"repeat_type"` // 每天1，工作日2，自定义3
	RepeatDate string            `json:"repeat_date"` // 自定义的情况下：1234567
}

func UpdateSilenceSetting(c *gin.Context) {
	var (
		req         SilenceSettingReq
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
	// 校验参数
	if err = req.validate(); err != nil {
		return
	}

	// 只有SA拥有者才能设置
	sessionUser = session.Get(c)
	err = req.UpdateNoDisturbingSetting(sessionUser.UserID, sessionUser.AreaID)
	if err != nil {
		return
	}

}

func (req *SilenceSettingReq) UpdateNoDisturbingSetting(userId int, areaID uint64) (err error) {

	setting := entity.MessageSilenceSetting{
		State:       req.State,
		IsAllDevice: req.IsAllDevice,
		EffectStart: time.Unix(req.EffectStart, 0),
		EffectEnd:   time.Unix(req.EffectEnd, 0),
		RepeatDate:  req.RepeatDate,
		RepeatType:  req.RepeatType,
	}
	err = entity.UpdateUserSetting(entity.SilenceSetting, &setting, userId, areaID)

	if err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}

	return nil
}

func (req *SilenceSettingReq) validate() (err error) {
	if !entity.CheckIllegalRepeatDate(req.RepeatDate) {
		err = errors.Wrap(err, status.MessageCenterRepeatDateIncorrectErr)
	}
	if req.RepeatType != entity.RepeatTypeAllDay && req.RepeatType != entity.RepeatTypeCustom {
		err = errors.Wrap(err, status.MessageCenterRepeatTypeIncorrectErr)
	}
	return
}
