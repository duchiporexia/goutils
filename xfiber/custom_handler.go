package xfiber

import (
	"github.com/duchiporexia/goutils/xerr"
	"github.com/gofiber/fiber/v2"
)

func CustomErrorHandler(ctx *fiber.Ctx, err error) error {
	if v, ok := err.(*xerr.BadRequestErr); ok {
		return SendBadRequestErr(ctx, v)
	}
	return SendInternalErr(ctx, err)
}
