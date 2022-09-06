package status

import "github.com/zhiting-tech/smartassistant/pkg/errors"

const (
	MessageCenterDeviceNoExist = iota + 13000
	MessageCenterTypeNoExist
	MessageCenterRepeatDateIncorrectErr
	MessageCenterRepeatTypeIncorrectErr
	MessageCenterTypeNotExistErr
)

func init() {
	errors.NewCode(MessageCenterDeviceNoExist, "设备不存在")
	errors.NewCode(MessageCenterTypeNoExist, "不存在的消息类型")
	errors.NewCode(MessageCenterRepeatDateIncorrectErr, "重复执行时间错误")
	errors.NewCode(MessageCenterRepeatTypeIncorrectErr, "重复执行配置错误")
	errors.NewCode(MessageCenterTypeNotExistErr, "消息类型不存在")
}
