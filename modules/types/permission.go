package types

import (
	"fmt"
)

const (
	FwUpgrade       = "firmware_upgrade" // 固件升级
	SoftwareUpgrade = "software_upgrade" // 软件升级
)

const (
	ActionAdd     = "add"
	ActionGet     = "get"
	ActionUpdate  = "update"
	ActionControl = "control"
	ActionDelete  = "delete"
	ActionManage  = "manage"
)

type Permission struct {
	ServiceName string `json:"service_name,omitempty"` // service名，仅针对设备属性权限
	Name        string `json:"name"`                   // attribute
	Action      string `json:"action"`                 // 动作
	Target      string `json:"target"`                 // 对象
	Attribute   string `json:"attribute"`              // 属性
}

// 设备
var (
	DeviceAdd         = Permission{Name: "添加设备", Action: ActionAdd, Target: "device"}
	DeviceUpdate      = Permission{Name: "修改设备", Action: ActionUpdate, Target: "device"}
	DeviceControl     = Permission{Name: "控制设备", Action: ActionControl, Target: "device"}
	DeviceDelete      = Permission{Name: "删除设备", Action: ActionDelete, Target: "device"}
	DeviceUpdateOrder = Permission{Name: "设备排序", Action: ActionUpdate, Target: "device", Attribute: "order"}
)

// 家庭/公司
var (
	AreaGetCode          = Permission{Name: "生成邀请码", Action: ActionGet, Target: "area", Attribute: "invite_code"}
	AreaUpdateName       = Permission{Name: "修改家庭名称", Action: ActionUpdate, Target: "area", Attribute: "name"}
	AreaUpdateMemberRole = Permission{Name: "修改成员角色", Action: ActionUpdate, Target: "area", Attribute: "member_role"}
	AreaDelMember        = Permission{Name: "删除成员", Action: ActionDelete, Target: "area", Attribute: "member"}
)

// 公司
var (
	AreaUpdateMemberDepartment = Permission{Name: "修改成员部门", Action: ActionUpdate, Target: "area", Attribute: "member_department"}
	AreaUpdateCompanyName      = Permission{Name: "修改公司名称", Action: ActionUpdate, Target: "area", Attribute: "company_name"}
)

// 房间/区域
var (
	LocationAdd         = Permission{Name: "添加房间/区域", Action: ActionAdd, Target: "location"}
	LocationUpdateOrder = Permission{Name: "调整顺序", Action: ActionUpdate, Target: "location", Attribute: "order"}
	LocationUpdateName  = Permission{Name: "修改房间名称", Action: ActionUpdate, Target: "location", Attribute: "name"}
	LocationGet         = Permission{Name: "查看房间详情", Action: ActionGet, Target: "location"}
	LocationDel         = Permission{Name: "删除房间", Action: ActionDelete, Target: "location"}
)

// 角色
var (
	RoleGet    = Permission{Name: "查看角色列表", Action: ActionGet, Target: "role"}
	RoleAdd    = Permission{Name: "新增角色", Action: ActionAdd, Target: "role"}
	RoleUpdate = Permission{Name: "编辑角色", Action: ActionUpdate, Target: "role"}
	RoleDel    = Permission{Name: "删除角色", Action: ActionDelete, Target: "role"}
)

// 场景
var (
	SceneAdd     = Permission{Name: "新增场景", Action: ActionAdd, Target: "scene"}
	SceneUpdate  = Permission{Name: "修改场景", Action: ActionUpdate, Target: "scene"}
	SceneDel     = Permission{Name: "删除场景", Action: ActionDelete, Target: "scene"}
	SceneControl = Permission{Name: "控制场景", Action: ActionControl, Target: "scene"}
)

// 部门
var (
	DepartmentAdd         = Permission{Name: "添加部门", Action: ActionAdd, Target: "department"}
	DepartmentUpdateOrder = Permission{Name: "调整部门顺序", Action: ActionUpdate, Target: "department", Attribute: "order"}
	DepartmentGet         = Permission{Name: "查看部门详情", Action: ActionGet, Target: "department"}
	DepartmentAddUser     = Permission{Name: "添加成员", Action: ActionAdd, Target: "department", Attribute: "user"}
	DepartmentUpdate      = Permission{Name: "部门设置", Action: ActionUpdate, Target: "department"}
)

var (
	DevicePermission     = []Permission{DeviceAdd, DeviceUpdate, DeviceControl, DeviceDelete, DeviceUpdateOrder}
	AreaPermission       = []Permission{AreaGetCode, AreaUpdateName, AreaUpdateMemberRole, AreaDelMember}
	LocationPermission   = []Permission{LocationAdd, LocationUpdateOrder, LocationUpdateName, LocationGet, LocationDel}
	RolePermission       = []Permission{RoleGet, RoleAdd, RoleUpdate, RoleDel}
	ScenePermission      = []Permission{SceneAdd, SceneUpdate, SceneDel, SceneControl}
	DepartmentPermission = []Permission{DepartmentAdd, DepartmentUpdateOrder, DepartmentGet, DepartmentAddUser, DepartmentUpdate}
	CompanyPermission    = []Permission{AreaGetCode, AreaUpdateCompanyName, AreaUpdateMemberRole, AreaUpdateMemberDepartment, AreaDelMember}
)

var (
	DefaultPermission []Permission
	// ManagerPermission 管理员默认权限
	ManagerPermission []Permission
	// MemberPermission 成员默认权限
	MemberPermission []Permission
)

func DeviceTarget(deviceID int) string {
	return fmt.Sprintf("device-%d", deviceID)
}

func NewDeviceDelete(deviceID int) Permission {
	target := DeviceTarget(deviceID)
	return Permission{Name: "删除设备", Action: ActionDelete, Target: target}
}
func NewDeviceUpdate(deviceID int) Permission {
	target := DeviceTarget(deviceID)
	return Permission{Name: "修改设备", Action: ActionUpdate, Target: target}
}

func NewDeviceManage(deviceID int, name string, attr string) Permission {
	target := DeviceTarget(deviceID)
	return Permission{Name: name, Action: ActionManage, Target: target, Attribute: attr}
}

func NewDeviceFwUpgrade(deviceID int) Permission {
	return NewDeviceManage(deviceID, "固件升级", FwUpgrade)
}

func NewDeviceSoftwareUpgrade(deviceID int) Permission {
	return NewDeviceManage(deviceID, "软件升级", SoftwareUpgrade)
}

func init() {

	DefaultPermission = append(DefaultPermission, DevicePermission...)
	DefaultPermission = append(DefaultPermission, AreaPermission...)
	DefaultPermission = append(DefaultPermission, LocationPermission...)
	DefaultPermission = append(DefaultPermission, RolePermission...)
	DefaultPermission = append(DefaultPermission, ScenePermission...)
	DefaultPermission = append(DefaultPermission, DepartmentPermission...)
	DefaultPermission = append(DefaultPermission, AreaUpdateMemberDepartment, AreaUpdateCompanyName)

	ManagerPermission = append(ManagerPermission, DefaultPermission...)
	MemberPermission = []Permission{DeviceControl, LocationGet, DepartmentGet}
}
