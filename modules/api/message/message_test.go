package message

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"testing"
	"time"
)

func GetDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("C:\\Users\\dev\\Desktop\\smartassistant\\data\\smartassistant\\sadb.db"), &gorm.Config{})
	if err != nil {
		logrus.Info("failed to connect database")
	}
	return db
}

func TestCreateUserSetting(t *testing.T) {

	db := GetDb()

	setting := entity.MessageDevNotificationSetting{
		Permission: true,
	}
	v, err := json.Marshal(setting)
	if err != nil {
		return
	}

	s := entity.UserSetting{
		Type:   "536s",
		UserId: 1,
		Value:  v,
		AreaId: 1,
	}
	err = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "type"}, {Name: "user_id"}, {Name: "area_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&s).Error
	fmt.Println(err)
}

func TestCreateMessageSetting(t *testing.T) {

	db := GetDb()
	s := entity.MessageSetting{
		DeviceId: 1,
		UserId:   1,
		State:    false,
		AreaId:   1,
	}
	msgSetting := entity.MessageSetting{}
	err := db.Table(s.TableName()).
		Where("user_id = ? and device_id = ? and area_id = ?", s.UserId, s.DeviceId, s.AreaId).
		First(&msgSetting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.Permission = true
			s.CreatedAt = time.Now()
			err = db.Create(&s).Error
		}
		return
	}
	err = db.Table(s.TableName()).Where("user_id = ? and device_id = ? and area_id = ?", s.UserId, s.DeviceId, s.AreaId).
		Update("state", s.State).Error
	fmt.Println(err)
}
