package xerr

import (
	"errors"
	"fmt"
	utilscommon "github.com/duchiporexia/goutils/internal"
	"github.com/duchiporexia/goutils/xlog"
)

type InternalErr struct {
	Code int
	Msg  string
}

func (e *InternalErr) Error() string {
	return e.Msg
}

func NewInternalErr(code int, msg string) error {
	return &InternalErr{
		Code: code,
		Msg:  msg,
	}
}

//////////////////////////////////////////////////////////
type BadRequestErr struct {
	Code int
	Msg  string
}

func (e *BadRequestErr) Error() string {
	return e.Msg
}
func NewBadRequestErr(msg string) error {
	return &BadRequestErr{
		Msg: msg,
	}
}
func NewBadRequestErrWithCode(code int, msg string) error {
	return &BadRequestErr{
		Code: code,
		Msg:  msg,
	}
}

type Str string

func (e Str) Error() string {
	return string(e)
}

func Is(err, target error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, target)
}

func LogAndHideErr(err error) error {
	if err == nil {
		return nil
	}
	msg := fmt.Sprintf("%+v", utilscommon.WithStack(err))
	xlog.ErrorE(2).Msg(msg)
	return ErrInternalErr
}
