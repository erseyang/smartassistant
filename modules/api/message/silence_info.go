package message

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
)

type SilenceSettingResp struct {
	State       bool `json:"state"`         // 免打扰状态，true 开启，false 关闭
	IsAllDevice bool `json:"is_all_device"` // 是否为所有设备，ture 是，false 否
	// 生效时间的配置
	EffectStart int64 `json:"effect_start"`
	EffectEnd   int64 `json:"effect_end"`

	// 重复执行的配置
	RepeatType int    `json:"repeat_type"` // 每天1，自定义2
	RepeatDate string `json:"repeat_date"` // 自定义的情况下：1234567
}

func SilenceSettingInfo(c *gin.Context) {
	var (
		err         error
		resp        SilenceSettingResp
		sessionUser = session.Get(c)
	)
	defer func() {
		response.HandleResponse(c, err, resp)
	}()

	setting := entity.GetDefaultSilenceMessageSetting()
	if err = entity.GetUserSetting(entity.SilenceSetting, &setting, sessionUser.AreaID); err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}

	resp.State = setting.State
	resp.IsAllDevice = setting.IsAllDevice
	resp.EffectStart = setting.EffectStart.Unix()
	resp.EffectEnd = setting.EffectEnd.Unix()
	resp.RepeatType = int(setting.RepeatType)
	resp.RepeatDate = setting.RepeatDate
}
