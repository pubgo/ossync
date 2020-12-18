package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pubgo/ossync/api/file"
	"github.com/pubgo/ossync/api/view"
	"github.com/pubgo/xerror"
)

func Router(r fiber.Router) {
	r.Use(func(view *fiber.Ctx) error {
		defer xerror.Resp(func(err xerror.XErr) {
			_ = view.JSON(fiber.Map{
				"code":   400,
				"detail": err,
				"msg":    err.Error(),
			})
		})

		return view.Next()
	})

	view.InitRouter(r.Group("/"))

	api := r.Group("/api")
	file.InitRouter(api)
}
