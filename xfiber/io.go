package xfiber

import (
	"github.com/duchiporexia/goutils/xerr"
	"github.com/duchiporexia/goutils/xlog"
	"github.com/gofiber/fiber/v2"
	"github.com/tinylib/msgp/msgp"
)

type msgToSend interface {
	EncodeMsg(en *msgp.Writer) (err error)
}
type msgToReceive interface {
	UnmarshalMsg(bts []byte) (o []byte, err error)
}
type msgWithCheck interface {
	Check() error
}

func Check(v msgWithCheck) error {
	return v.Check()
}

func ReceiveMsgAndCheck(ctx *fiber.Ctx, v interface {
	msgWithCheck
	msgToReceive
}) error {
	err := ReceiveMsg(ctx, v)
	if err == nil {
		err = Check(v)
	}
	return err
}

func ReceiveMsg(ctx *fiber.Ctx, v msgToReceive) error {
	left, err := v.UnmarshalMsg(ctx.Body())
	if err != nil {
		return err
	}
	if len(left) > 0 {
		xlog.Warn("left > 0")
	}
	return nil
}

func SendMsg(ctx *fiber.Ctx, v msgToSend) error {
	writer := msgp.NewWriter(ctx.Response().BodyWriter())
	err := v.EncodeMsg(writer)
	if err != nil {
		xlog.ErrorE(1).Err(err).Send()
	}
	return writer.Flush()
}

func SendErr(ctx *fiber.Ctx, status int, err error) error {
	ctx.Status(status)
	eErr, ok := err.(*xerr.InternalErr)
	var oerr error
	if ok {
		oerr = SendMsg(ctx, ErrResponse{Code: eErr.Code, Msg: eErr.Msg})
	} else {
		oerr = SendMsg(ctx, ErrResponse{Code: status, Msg: err.Error()})
	}
	if oerr != nil {
		xlog.Error(oerr)
	}
	return nil
}

func SendInternalErr(ctx *fiber.Ctx, err error) error {
	return SendErr(ctx, fiber.StatusInternalServerError, err)
}

func SendBadRequestErr(ctx *fiber.Ctx, err error) error {
	ctx.Status(fiber.StatusBadRequest)
	eErr, ok := err.(*xerr.BadRequestErr)
	var oerr error
	if ok {
		oerr = SendMsg(ctx, ErrResponse{Code: eErr.Code, Msg: eErr.Msg})
	} else {
		oerr = SendMsg(ctx, ErrResponse{Code: fiber.StatusBadRequest, Msg: err.Error()})
	}
	if oerr != nil {
		xlog.Error(oerr)
	}
	return nil
}
