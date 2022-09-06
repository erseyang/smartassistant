package message

const (
	TypePlatformAuth = iota + 1
)

const (
	ContentMemberJoin       = "成员（%v）通过%v分享的二维码加入%v"
	ContentMemberExit       = "成员（%v）已退出%v"
	ContentDeviceAdd        = "成员（%v）已添加设备（%v）"
	ContentDeviceDel        = "成员（%v）已删除设备（%v）"
	ContentPlatformAuth     = "成员（%v）已授权第三方平台（%v）"
	ContentPlatformUnbind   = "成员（%v）已解除授权第三方平台（%v）"
	ContentHomeBridgeAuth   = "成员（%v）已授权家居桥接"
	ContentHomeBridgeUnbind = "成员（%v）已解除授权家居桥接"
)

type NotificationMessage struct {
	ID          int    `json:"id"`
	ReceiverId  int    `json:"receiver_id"`
	Type        int    `json:"type"` // 1 告警 ，2 家庭/公司
	Location    string `json:"location"`
	CreatedTime int64  `json:"created_time"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	LogoUrl     string `json:"logo_url"`
}