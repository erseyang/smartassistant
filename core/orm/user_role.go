package orm

import (
	"gitlab.yctc.tech/root/smartassistent.git/utils/errors"
	"gorm.io/gorm"
)

type UserRole struct {
	ID     int
	UserID int `gorm:"unique_index:uid_rid"`
	RoleID int `gorm:"unique_index:uid_rid"`
}

func (ur UserRole) TableName() string {
	return "user_roles"
}

func CreateUserRole(uRoles []UserRole) (err error) {
	if err = GetDB().Create(&uRoles).Error; err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}
	return
}

func GetRoleIdsByUid(userId int) (roleIds []int, err error) {

	if err = GetDB().Model(&UserRole{}).Where("user_id = ?", userId).Pluck("role_id", &roleIds).Error; err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}
	return
}

func GetRolesByUid(userId int) (roles []Role, err error) {
	if err = GetDB().Model(&Role{}).
		Joins("inner join user_roles on roles.id=user_roles.role_id").
		Where("user_roles.user_id = ?", userId).Find(&roles).Error; err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}
	return
}

func DelUserRoleByUid(userId int, db *gorm.DB) (err error) {
	err = db.Where("user_id=?", userId).Delete(&UserRole{}).Error
	if err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
	}
	return
}