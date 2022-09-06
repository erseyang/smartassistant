package message

import (
	"encoding/json"
	"fmt"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/plugin"
	"github.com/zhiting-tech/smartassistant/modules/types"
	"github.com/zhiting-tech/smartassistant/modules/utils/url"
	"github.com/zhiting-tech/smartassistant/pkg/event"
	"github.com/zhiting-tech/smartassistant/pkg/logger"
	"strconv"
	"sync"
	"time"
)

var messagesManager *MessagesManager
var messageOnce sync.Once

func GetMessagesManager() *MessagesManager {
	messageOnce.Do(func() {
		messagesManager = newMessagesManager()
	})
	return messagesManager
}

func newMessagesManager() *MessagesManager {
	return &MessagesManager{
		messagesCh: make(chan entity.MessageRecord, 128),
	}
}

type MessagesManager struct {
	messagesCh chan entity.MessageRecord
}

func (m *MessagesManager) SendMsg(msg entity.MessageRecord) {
	m.messagesCh <- msg
}

func (m *MessagesManager) Run() {
	for {
		select {
		case msg := <-m.messagesCh:
			m.saveAndSendMsg(msg)
		}
	}
}

func (m *MessagesManager) saveAndSendMsg(msg entity.MessageRecord) {
	var (
		uIds       []int
		err        error
		usList     []entity.UserSetting
		msList     []entity.MessageSetting
		msgRecords []entity.MessageRecord
	)
	if uIds, err = entity.GetUIds(msg.AreaId); err != nil {
		logger.Debug("SaveAndSendMsg GetUIds error", err)
		return
	}

	if msg.Type == entity.NotifyTypeAlarm {
		if usList, err = entity.GetUserSettingsByAreaId(msg.AreaId); err != nil {
			logger.Warning("SaveAndSendMsg GetUserSettingsByAreaId error", err)
			return
		}

		if msList, err = entity.GetMessageSettingsByAreaId(msg.AreaId); err != nil {
			logger.Warning("SaveAndSendMsg GetMessageSettingsByAreaId error", err)
			return
		}
	}
	for _, id := range uIds {

		if FilterMsg(id, msg.DeviceId, msg.CreatedAt, usList, msList) {
			continue
		}

		msgRecord := entity.MessageRecord{
			Type:       msg.Type,
			ReceiverId: id,
			AreaId:     msg.AreaId,
			LocationId: msg.LocationId,
			DeviceId:   msg.DeviceId,
			Title:      msg.Title,
			Content:    msg.Content,
			CreatedAt:  msg.CreatedAt,
		}
		msgRecords = append(msgRecords, msgRecord)
	}

	msgRecordList, err := entity.CreateMessageRecords(msgRecords)
	if err != nil {
		logger.Debug("Create Message Records error", err)
		return
	}

	m.SendSameMsg(msgRecordList)
}

func (m *MessagesManager) SendSameMsg(msgRecords []entity.MessageRecord) {
	var (
		err    error
		dev    entity.Device
		areaId uint64
	)
	if len(msgRecords) > 0 {
		if msgRecords[0].DeviceId != 0 && msgRecords[0].Type == entity.NotifyTypeAlarm {
			if dev, err = entity.GetDeviceByID(msgRecords[0].DeviceId); err != nil {
				logger.Warning("MsgSliRecordToNotification GetDeviceByID error ", err.Error())
			}
		}
		areaId = msgRecords[0].AreaId
	}

	for _, msg := range MsgRecordsToNotifications(msgRecords) {
		em := event.NewEventMessage(event.MessageCenter, areaId)
		if msg.Type == entity.NotifyTypeAlarm {
			if dev.LogoType == nil || *dev.LogoType == int(types.NormalLogo) {
				logo := plugin.GetGlobalClient().Config(dev.PluginID).DeviceConfig(dev.Model, dev.Type).Logo
				msg.LogoUrl = plugin.ConcatPluginPath(dev.PluginID, logo)
			} else {
				if logoInfo, ok := types.GetLogo(types.LogoType(*dev.LogoType)); ok {
					msg.LogoUrl = url.ImagePath(logoInfo.FileName)
				}
			}
		} else if msg.Type == entity.NotifyTypeArea {
			msg.LogoUrl = url.ImagePath(entity.AreaLogoUrlMap[msg.Title])
		}
		em.Param["message"] = msg
		em.Param["user_id"] = msg.ReceiverId
		event.Notify(em)
	}
	return
}

// FilterMsg 判断是否过滤消息
func FilterMsg(uid, devId int, createdAt time.Time, usList []entity.UserSetting, msList []entity.MessageSetting) bool {

	notifyFlag := false
	noDisturbingFlag := false
	for _, us := range usList {
		// 没有个人设置时，默认入库发送
		if uid == us.UserId {
			if us.Type == entity.DevNotificationSetting { // 设备通知设置
				setting := entity.GetDefaultDevNotificationMessageSetting()
				_ = json.Unmarshal(us.Value, &setting)
				// 设备通知设置为false时，都不入库推送
				if setting.Permission == false {
					notifyFlag = true
					continue
				} else {
					// 设备通知设置为ture时，查询消息设置表是否有静止该设备通知，有的话也不入库推送
					for _, ms := range msList {
						if uid == ms.UserId && devId == ms.DeviceId {
							if ms.Permission == false {
								notifyFlag = true
								continue
							}
						}
					}
				}
			} else if us.Type == entity.SilenceSetting { // 个人免打扰设置
				setting := entity.GetDefaultSilenceMessageSetting()
				_ = json.Unmarshal(us.Value, &setting)
				if setting.State == true && IsInTimePeriod(createdAt, setting.EffectStart, setting.EffectEnd) &&
					(setting.RepeatType == entity.RepeatTypeAllDay || IsInDay(createdAt, setting.RepeatDate)) {
					if setting.IsAllDevice == true {
						noDisturbingFlag = true
					} else {
						notifyFlag = true
						for _, ms := range msList {
							if uid == ms.UserId && devId == ms.DeviceId {
								if ms.State == false {
									notifyFlag = false
									continue
								}
							}
						}
					}

				}
			}
		}
	}

	if notifyFlag {
		return true
	} else if noDisturbingFlag {
		return true
	}
	return false
}

// IsInTimePeriod 是否在时间段内
func IsInTimePeriod(t, effectStart, effectEnd time.Time) bool {

	days := int(t.Sub(effectStart).Hours() / 24)
	effectEndTime := effectEnd.AddDate(0, 0, days)
	effectStartTime := effectStart.AddDate(0, 0, days)
	return t.Before(effectEndTime) && time.Now().After(effectStartTime)

}

func IsInDay(day time.Time, repeatDate string) bool {
	weekDay := strconv.Itoa(int(day.Weekday()))
	for _, date := range repeatDate {
		if fmt.Sprintf("%c", date) == weekDay {
			return true
		}
	}
	return false
}

func MsgRecordsToNotifications(msgRecords []entity.MessageRecord) (msgList []NotificationMessage) {
	var (
		err      error
		location entity.Location
	)
	if len(msgRecords) > 0 {
		if msgRecords[0].LocationId != 0 {
			if location, err = entity.GetLocationByID(msgRecords[0].LocationId); err != nil {
				logger.Warningf("MsgSliRecordToNotification GetLocationByID error", err)
			}
		}
	}
	for _, msgRecord := range msgRecords {
		msgList = append(msgList, NotificationMessage{
			ID:          msgRecord.ID,
			ReceiverId:  msgRecord.ReceiverId,
			Type:        msgRecord.Type,
			Location:    location.Name,
			CreatedTime: msgRecord.CreatedAt.Unix(),
			Title:       msgRecord.Title,
			Content:     msgRecord.Content,
		})
	}
	return
}
