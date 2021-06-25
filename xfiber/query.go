package xfiber

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func Query(ctx *fiber.Ctx, key string, defaultValue string, rules ...validation.Rule) (string, error) {
	str := ctx.Query(key)
	if str == "" {
		return defaultValue, nil
	}
	if len(rules) > 0 {
		err := validation.Validate(str, rules...)
		if err != nil {
			return "", err
		}
	}
	return str, nil
}

func QueryInt(ctx *fiber.Ctx, key string, defaultValue int, rules ...validation.Rule) (int, error) {
	str := ctx.Query(key)
	if str == "" {
		return defaultValue, nil
	}
	value, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	if err != nil {
		return 0, err
	}
	if len(rules) > 0 {
		err = validation.Validate(value, rules...)
		if err != nil {
			return 0, err
		}
	}
	return value, nil
}
