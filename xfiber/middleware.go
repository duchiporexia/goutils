package xfiber

import (
	"github.com/duchiporexia/goutils/xgenid"
	"github.com/duchiporexia/goutils/xlog"
	"github.com/gofiber/fiber/v2"
)

func NewAPILogger() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		xlog.InfoE().Str("method", c.Method()).Str("reqid", string(c.Response().Header.Peek(fiber.HeaderXRequestID))).Msg(c.OriginalURL())
		return c.Next()
	}
}

func NewRequestId() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Get id from request
		rid := c.Get(fiber.HeaderXRequestID)
		// Create new id if empty
		if rid == "" {
			rid = xgenid.Uuid()
		}
		// Set new id to response
		c.Set(fiber.HeaderXRequestID, rid)
		// Bye
		return c.Next()
	}
}
