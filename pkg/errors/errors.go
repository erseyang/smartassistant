package errors

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"

	"github.com/pkg/errors"
)

type Error struct {
	Err  error
	Code Code
}

func (e Error) Error() string {
	return e.Code.Reason
}

func (e Error) GRPCStatus() *status.Status {
	return status.New(codes.Code(e.Code.Status), e.Code.Reason)
}

// func (e Error) Format(f fmt.State, verb rune) {
// 	io.WriteString(f, e.Error())
// 	stackTrace := e.GetErrStack()
// 	stackTrace.Format(f, verb)
// }

func (e Error) GetErrStack() errors.StackTrace {
	// 获取错误调用栈(跳过New,Newf,Wrapf,Wrap调用栈)
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	st := e.Err.(stackTracer)
	stackTrace := st.StackTrace()

	filterStack := errors.StackTrace{}
	filterFuncRegex, _ := regexp.Compile(`/utils/errors\.(New|Wrap)f?`)

	for _, f := range stackTrace[:2] {
		stackText, _ := f.MarshalText()
		if !filterFuncRegex.MatchString(string(stackText)) {
			filterStack = append(filterStack, f)
		}
	}
	filterStack = append(filterStack, stackTrace[2:]...)
	return filterStack

}

// New
func New(status int) error {
	return Newf(status)
}

// Newf
func Newf(status int, args ...interface{}) error {

	code := GetCode(status)
	if len(args) != 0 {
		code.Reason = fmt.Sprintf(code.Reason, args...)
	}
	return Error{
		Err:  errors.New(code.Reason),
		Code: code,
	}
}

// Wrap
func Wrap(err error, status int) error {
	return Wrapf(err, status)
}

// WrapErrorf 支持格式化错误原因
func Wrapf(err error, status int, args ...interface{}) error {

	code := GetCode(status)
	switch v := err.(type) {
	case Error:
		err = v.Err
		code = v.Code
	default:
		err = errors.WithStack(err)
	}
	if len(args) != 0 {
		code.Reason = fmt.Sprintf(code.Reason, args...)
	}

	if err == nil { // 避免传nil error时日志没有任何调用栈
		err = errors.New(code.Reason)
	}

	return Error{
		Err:  err,
		Code: code,
	}
}

// Cause 获取原始错误
func Cause(err error) error {
	switch v := err.(type) {
	case Error:
		return errors.Cause(v.Err)
	default:
		return errors.Cause(err)
	}
}
