package entity

import (
	"encoding/json"
	"errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// 默认配置表
var (
	defaultUserSettingMap = map[string]interface{}{
		SilenceSetting:         defaultSilenceSetting,
		DevNotificationSetting: defaultDevNotificationSetting,
		MsgCenterReadMsgId:     defaultMessageReadMsgId,
	}
)

// 配置类型
const (
	SilenceSetting         = "silence_setting"
	DevNotificationSetting = "device_notification_setting"
	MsgCenterReadMsgId     = "msg_center_read_msg_id"
)

// 默认配置项
var (
	defaultSilenceSetting = MessageSilenceSetting{
		State:       false,
		IsAllDevice: true,
		RepeatType:  RepeatTypeAllDay,
		EffectStart: time.Date(2020, 1, 1, 23, 0, 0, 0, time.Local),
		EffectEnd:   time.Date(2020, 1, 2, 7, 0, 0, 0, time.Local),
	}
	defaultDevNotificationSetting = MessageDevNotificationSetting{
		Permission: true,
	}

	defaultMessageReadMsgId = MessageReadMsgId{
		MsgId: 0,
	}
)

// MessageSilenceSetting 用户消息免打扰配置
type MessageSilenceSetting struct {
	State       bool `json:"state"`         // 免打扰状态，true 开启，false 关闭
	IsAllDevice bool `json:"is_all_device"` // 是否为所有设备，ture 是，false 否
	// 生效时间的配置
	EffectStart time.Time `json:"effect_start"`
	EffectEnd   time.Time `json:"effect_end"`

	// 重复执行的配置
	RepeatType RepeatType `json:"repeat_type"` // 每天1，工作日2，自定义3
	RepeatDate string     `json:"repeat_date"` // 自定义的情况下：1234567
}

// MessageDevNotificationSetting 用户消息设备通知配置
type MessageDevNotificationSetting struct {
	Permission bool `json:"permission"`
}

// MessageReadMsgId 用户已读消息最大id
type MessageReadMsgId struct {
	MsgId int `json:"msg_id"`
}

// UserSetting 用户个人设置表
type UserSetting struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	AreaId    uint64 `json:"area_id" gorm:"uniqueIndex:user_area_type,priority:2"`
	UserId    int    `json:"user_id" gorm:"uniqueIndex:user_area_type,priority:1"`
	Type      string `json:"type" gorm:"uniqueIndex:user_area_type,priority:3"`
	Value     datatypes.JSON
	CreatedAt time.Time `json:"created_at"`
}

func (m UserSetting) TableName() string {
	return "user_setting"
}

func CreateUserSetting(setting UserSetting) (id int, err error) {
	err = GetDB().Create(&setting).Error
	id = setting.UserId
	return
}

// GetUserSetting 获取个人设置
func GetUserSetting(settingType string, setting interface{}, areaID uint64) (err error) {
	var us UserSetting
	err = GetDB().Where(UserSetting{Type: settingType, AreaId: areaID}).First(&us).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}

	return json.Unmarshal(us.Value, setting)
}

func GetUserSettingsByAreaId(areaID uint64) (usList []UserSetting, err error) {
	err = GetDB().Where(UserSetting{AreaId: areaID}).Find(&usList).Error
	return
}

// UpdateUserSetting 添加用户设置
func UpdateUserSetting(settingType string, setting interface{}, userID int, areaId uint64) (err error) {

	v, err := json.Marshal(setting)
	if err != nil {
		return
	}

	s := UserSetting{
		Type:   settingType,
		UserId: userID,
		Value:  v,
		AreaId: areaId,
	}
	return GetDB().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "type"}, {Name: "user_id"}, {Name: "area_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&s).Error
}

// GetDefaultSilenceMessageSetting 获取用户消息免打扰默认配置
func GetDefaultSilenceMessageSetting() MessageSilenceSetting {
	return defaultUserSettingMap[SilenceSetting].(MessageSilenceSetting)
}

// GetDefaultDevNotificationMessageSetting 获取用户消息设备通知默认配置
func GetDefaultDevNotificationMessageSetting() MessageDevNotificationSetting {
	return defaultUserSettingMap[DevNotificationSetting].(MessageDevNotificationSetting)
}

// GetDefaultMessageReadMsgId 获取用户已读消息默认id
func GetDefaultMessageReadMsgId() MessageReadMsgId {
	return defaultUserSettingMap[MsgCenterReadMsgId].(MessageReadMsgId)
}