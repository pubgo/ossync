package view

import (
	"github.com/gofiber/fiber/v2"
)

func InitRouter(r fiber.Router) {
	r.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", nil)
	})
	r.Get("/task", func(ctx *fiber.Ctx) error {
		return ctx.Render("task", nil)
	})
}
