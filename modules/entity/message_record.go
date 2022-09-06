package entity

import (
	"time"
)

const (
	TitleMemberChange     = "成员变更"
	TitleDeviceChange     = "设备变更"
	TitlePlatformChange   = "第三方平台变更"
	TitleHomeBridgeChange = "家居桥接变更"
)

var AreaLogoUrlMap = map[string]string{
	TitleMemberChange:     "member_change.png",
	TitleDeviceChange:     "device_change.png",
	TitlePlatformChange:   "platform_change.png",
	TitleHomeBridgeChange: "homebridge_change.png",
}

// 通知类型
const (
	NotifyTypeAlarm = iota + 1
	NotifyTypeArea
)

// MessageRecord 消息记录表
type MessageRecord struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	AreaId     uint64    `json:"area_id"`
	DeviceId   int       `json:"device_id"`
	LocationId int       `json:"location_id"`
	Type       int       `json:"type"`
	ReceiverId int       `json:"receiver_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

func (m MessageRecord) TableName() string {
	return "message_record"
}

func NewMessageRecord(areaId uint64, title, content string) MessageRecord {
	return MessageRecord{
		AreaId:    areaId,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
	}
}

func NewAlarmMessageRecord(areaId uint64, deviceId, locationId int, title, content string) (msgRecord MessageRecord) {
	msgRecord = NewMessageRecord(areaId, title, content)
	msgRecord.Type = NotifyTypeAlarm
	msgRecord.DeviceId = deviceId
	msgRecord.LocationId = locationId
	return
}

func NewAreaMessageRecord(areaId uint64, title, content string) (msgRecord MessageRecord) {
	msgRecord = NewMessageRecord(areaId, title, content)
	msgRecord.Type = NotifyTypeArea
	return
}

func NewMemberChangeMessageRecord(areaId uint64, content string) (msgRecord MessageRecord) {
	msgRecord = NewAreaMessageRecord(areaId, TitleMemberChange, content)
	return
}

func NewHomeBridgeChangeMessageRecord(areaId uint64, content string) (msgRecord MessageRecord) {
	msgRecord = NewAreaMessageRecord(areaId, TitleHomeBridgeChange, content)
	return
}

func NewDeviceChangeMessageRecord(areaId uint64, content string) (msgRecord MessageRecord) {
	msgRecord = NewAreaMessageRecord(areaId, TitleDeviceChange, content)
	return
}

func NewPlatformChangeMessageRecord(areaId uint64, content string) (msgRecord MessageRecord) {
	msgRecord = NewAreaMessageRecord(areaId, TitlePlatformChange, content)
	return
}

func CreateMessageRecords(msgRecords []MessageRecord) (msgList []MessageRecord, err error) {

	err = GetDB().Create(&msgRecords).Error
	if err != nil {
		return
	}

	msgList = append(msgRecords)

	return
}

type MessageRecordInfo struct {
	ID           int       `json:"id" `
	LocationName string    `json:"location_name"`
	Type         int       `json:"type"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	PluginId     string    `json:"plugin_id"`
	Model        string    `json:"model"`
	DeviceType   string    `json:"device_type"`
	CreatedAt    time.Time `json:"created_at"`
}

func MessageRecordInfoList(mid, uid, count, msgType int, areaId uint64) (msgRecords []MessageRecordInfo, err error) {

	fields := []string{"message_record.id as id", "message_record.type as type", "message_record.title",
		"message_record.content as content", "message_record.created_at as created_at", "devices.plugin_id",
		"devices.model as model", "devices.type as device_type", "locations.name as location_name"}

	if msgType == 0 {
		err = GetDB().Table(MessageRecord{}.TableName()).Select(fields).
			Joins("left join " + Device{}.TableName() + " on message_record.device_id = devices.id").
			Joins("left join " + Location{}.TableName() + " on message_record.location_id = locations.id").
			Where("message_record.id < ? and message_record.receiver_id = ? and " +
				"message_record.area_id = ?", mid, uid, areaId).
			Order("id desc").Limit(count).Find(&msgRecords).Error
	} else {
		err = GetDB().Table(MessageRecord{}.TableName()).Select(fields).
			Joins("left join " + Device{}.TableName() + " on message_record.device_id = devices.id").
			Joins("left join " + Location{}.TableName() + " on message_record.location_id = locations.id").
			Where("message_record.id < ? and message_record.receiver_id = ? and "+
				"message_record.type = ? and message_record.area_id = ?", mid, uid, msgType, areaId).
			Order("id desc").Limit(count).Find(&msgRecords).Error
	}

	return
}

func MessageRecordLatestList(uid, count, msgType int, areaId uint64) (msgRecords []MessageRecordInfo, err error) {
	fields := []string{"message_record.id as id", "message_record.type as type", "message_record.title",
		"message_record.content as content", "message_record.created_at as created_at", "devices.plugin_id",
		"devices.model as model", "devices.type as device_type", "locations.name as location_name"}


	if msgType == 0 {
		err = GetDB().Table(MessageRecord{}.TableName()).Select(fields).
			Joins("left join " + Device{}.TableName() + " on message_record.device_id = devices.id").
			Joins("left join " + Location{}.TableName() + " on message_record.location_id = locations.id").
			Where("message_record.receiver_id = ? and message_record.area_id = ?", uid, areaId).
			Order("id desc").Limit(count).Find(&msgRecords).Error
	} else {
		err = GetDB().Table(MessageRecord{}.TableName()).Select(fields).
			Where("message_record.receiver_id = ? and message_record.type = ? and message_record.area_id = ?", uid, msgType, areaId).
			Joins("left join " + Device{}.TableName() + " on message_record.device_id = devices.id").
			Joins("left join " + Location{}.TableName() + " on message_record.location_id = locations.id").
			Order("id desc").Limit(count).Find(&msgRecords).Error
	}
	return
}

func GetLatestMsgRecord(uid int, areaId uint64) (msgRecord MessageRecord, err error) {
	err = GetDB().Table(MessageRecord{}.TableName()).
		Where("receiver_id = ? and area_id = ?", uid, areaId).
		Order("id desc").
		First(&msgRecord).Error
	return
}
