package types

const (
	SATokenKey      = "smart-assistant-token"
	SaModel         = "MH-SA00DCLB001W"
	VerificationKey = "verification-code" // 临时密码

	RoleKey   = "role"
	OwnerRole = "owner"

	// SAID 云端校验来自sa的请求时使用
	SAID  = "SA-ID"
	SAKey = "SA-Key"

	// GrantType 授权方式
	GrantType = "Grant-Type"

	SAMaintenanceTokenKey = "sa-maintenance-token"
)

const (
	// NormalMode 运行模式:正常模式
	NormalMode = iota

	// MaintenanceModeDiscover 运行模式:维护模式等待连接
	MaintenanceModeDiscover

	// MaintenanceModeConnected 运行模式:维护模式已连接
	MaintenanceModeConnected

	// MaintenanceModeResetPassword 运行模式:维护模式已连接专业版重置重置密码无需旧密码
	MaintenanceModeResetPassword
)

const (
	CloudDisk     = "wangpan"
	CloudDiskAddr = "wangpan:8089"

	HomeBridge = "homebridge"
)
