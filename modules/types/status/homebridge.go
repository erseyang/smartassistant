package status

import "github.com/zhiting-tech/smartassistant/pkg/errors"

const (
	GeneratePinError = iota + 12000
	UnbindError
)

func init() {
	errors.NewCode(GeneratePinError, "pin 码生成失败")
	errors.NewCode(UnbindError, "解除授权失败")
}
