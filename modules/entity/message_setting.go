package entity

import (
	"gorm.io/gorm"
	"time"
)

// MessageSetting 消息设置表
type MessageSetting struct {
	ID         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	AreaId     uint64 `json:"area_id"`
	DeviceId   int    `json:"device_id"`
	UserId     int    `json:"user_id"`
	Permission bool   `json:"permission"` // true 允许通知， false 禁止通知
	State      bool   `json:"state"`      // true 免打扰, false 关闭免打扰

	CreatedAt time.Time `json:"created_at"`
}

func (m MessageSetting) TableName() string {
	return "message_setting"
}

// UpdateMessageSettingState 更新||添加
func UpdateMessageSettingState(setting MessageSetting) (err error) {
	msgSetting := MessageSetting{}
	err = GetDB().Table(setting.TableName()).
		Where("user_id = ? and device_id = ? and area_id = ?", setting.UserId, setting.DeviceId, setting.AreaId).
		First(&msgSetting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			setting.Permission = true
			setting.CreatedAt = time.Now()
			err = GetDB().Create(&setting).Error
		}
		return
	}
	err = GetDB().Table(setting.TableName()).
		Where("user_id = ? and device_id = ? and area_id = ?", setting.UserId, setting.DeviceId, setting.AreaId).
		Update("state", setting.State).Error
	return
}

func UpdateMessageSettingPermission(setting MessageSetting) (err error) {
	msgSetting := MessageSetting{}
	err = GetDB().Table(setting.TableName()).
		Where("user_id = ? and device_id = ? and area_id = ?", setting.UserId, setting.DeviceId, setting.AreaId).
		First(&msgSetting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			setting.State = true
			setting.CreatedAt = time.Now()
			err = GetDB().Create(&setting).Error
		}
		return
	}
	err = GetDB().Table(setting.TableName()).
		Where("user_id = ? and device_id = ? and area_id = ?", setting.UserId, setting.DeviceId, setting.AreaId).
		Update("permission", setting.Permission).Error
	return
}

func GetMessageSettings(areaId uint64, userId int) (msgSettings []MessageSetting, err error) {
	err = GetDB().Table(MessageSetting{}.TableName()).
		Where("area_id = ? and user_id = ?", areaId, userId).
		Find(&msgSettings).Error
	return
}

func GetMessageSettingsByAreaId(areaId uint64) (msgSettings []MessageSetting, err error) {
	err = GetDB().Table(MessageSetting{}.TableName()).
		Where("area_id = ?", areaId).
		Find(&msgSettings).Error
	return
}

func GetMessageSetting(areaId uint64, userId, deviceId int) (msgSetting MessageSetting, err error) {
	err = GetDB().Table(MessageSetting{}.TableName()).
		Where("area_id = ? and user_id = ? and device_id = ?", areaId, userId, deviceId).
		First(&msgSetting).Error
	return
}
