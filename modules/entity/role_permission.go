package entity

import (
	"strconv"
	"strings"

	"gorm.io/gorm"

	"github.com/zhiting-tech/smartassistant/modules/types"
)

type ActionType string

// RolePermission 角色权限role_tmpl
type RolePermission struct {
	ID        int
	RoleID    int  `gorm:"index:permission,unique"` // 角色
	Role      Role `gorm:"constraint:OnDelete:CASCADE;"`
	Name      string
	Action    string `gorm:"index:permission,unique"` // 动作
	Target    string `gorm:"index:permission,unique"` // 对象
	Attribute string `gorm:"index:permission,unique"` // 属性
}

func (p RolePermission) TableName() string {
	return "role_permissions"
}

// IsDeviceControlPermit 判断用户是否有该设备的某个控制权限
func IsDeviceControlPermit(userID, deviceID int, attr Attribute) bool {
	return IsDeviceControlPermitByAttr(userID, deviceID, attr.Attribute.AID)
}

// IsDeviceControlPermitByAttr 判断用户是否有该设备的某个控制权限
func IsDeviceControlPermitByAttr(userID, deviceID int, aid int) bool {
	target := types.DeviceTarget(deviceID)
	return judgePermit(userID, types.ActionControl, target, strconv.Itoa(aid))
}

type Attr struct {
	DeviceID   int
	InstanceID int
	Attribute  string
}

type Permissions struct {
	ps       []RolePermission
	isOwner  bool
	userID   int
	clientID string
	scope    string
}

func (up Permissions) IsClient() bool {
	return up.clientID != "" && up.userID == 0
}

// DeviceScopeAllow 是否允许控制设备
func (up Permissions) DeviceScopeAllow() bool {
	return strings.Contains(up.scope, types.ScopeDevice.Scope) ||
		strings.Contains(up.scope, types.ScopeGetTokenBySC.Scope) // 兼容旧的SC Client
}

func (up Permissions) IsOwner() bool {
	return up.isOwner
}

// IsDeviceControlPermit 判断设备是否可控制
func (up Permissions) IsDeviceControlPermit(deviceID int) bool {
	if up.isOwner {
		return true
	}
	for _, p := range up.ps {
		if p.Action == types.ActionControl &&
			p.Target == types.DeviceTarget(deviceID) {
			return true
		}
	}
	return false
}

// IsDeviceAttrControlPermit 判断设备的属性是否有权限
func (up Permissions) IsDeviceAttrControlPermit(deviceID int, aid int) bool {
	if up.isOwner {
		return true
	}
	if up.IsClient() {
		return up.DeviceScopeAllow() // 有设备控制权限的client
	}
	for _, p := range up.ps {
		if p.Action == types.ActionControl &&
			p.Target == types.DeviceTarget(deviceID) &&
			p.Attribute == strconv.Itoa(aid) {
			return true
		}
	}
	return false
}

func (up Permissions) IsPermit(tp types.Permission) bool {
	if up.isOwner {
		return true
	}
	for _, p := range up.ps {
		if p.Action == tp.Action && p.Target == tp.Target && p.Attribute == tp.Attribute {
			return true
		}
	}
	return false
}

// GetUserPermissions 获取用户的所有权限
func GetUserPermissions(userID int) (up Permissions, err error) {
	var ps []RolePermission
	if err = GetDB().Scopes(UserRolePermissionsScope(userID)).
		Find(&ps).Error; err != nil {
		return
	}
	user, err := GetUserByID(userID)
	if err != nil {
		return
	}
	up = Permissions{
		userID: userID, ps: ps,
		isOwner: IsOwnerOfArea(userID, user.AreaID)}
	return
}

// GetClientPermissions 获取客户的权限
func GetClientPermissions(clientID, scope string) (up Permissions, err error) {
	up = Permissions{
		clientID: clientID,
		scope:    scope,
	}
	return
}

func UserRolePermissionsScope(userID int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select("role_permissions.*").
			Joins("inner join roles on roles.id=role_permissions.role_id").
			Joins("inner join user_roles on user_roles.role_id=roles.id").
			Where("user_roles.user_id=?", userID)
	}
}
func JudgePermit(userID int, p types.Permission) bool {
	return judgePermit(userID, p.Action, p.Target, p.Attribute)
}

func judgePermit(userID int, action, target, attribute string) bool {
	// SA拥有者默认拥有所有权限
	user, err := GetUserByID(userID)
	if err != nil {
		return false
	}
	if IsOwnerOfArea(userID, user.AreaID) {
		return true
	}

	var permissions []RolePermission
	if err := GetDB().Scopes(UserRolePermissionsScope(userID)).
		Where("action = ? and target = ? and attribute = ?",
			action, target, attribute).Find(&permissions).Error; err != nil {
		return false
	}

	if len(permissions) == 0 {
		return false
	}

	return true
}

func IsPermit(roleID int, action, target, attribute string, tx *gorm.DB) bool {
	p := RolePermission{
		RoleID:    roleID,
		Action:    action,
		Target:    target,
		Attribute: attribute,
	}
	if err := tx.First(&p, p).Error; err != nil {
		return false
	}
	return true
}

func IsDeviceActionPermit(roleID int, action string, tx *gorm.DB) bool {
	return IsPermit(roleID, action, "device", "", tx)
}
