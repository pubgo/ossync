package file

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pubgo/ossync/models"
	"github.com/pubgo/xerror"
)

func InitRouter(r fiber.Router) {
	r.Delete("/file/:id", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		models.SyncFileDelete("id=?", id)
		return nil
	})
	r.Put("/file/:id", func(ctx *fiber.Ctx) error {
		var sf map[string]interface{}
		xerror.Panic(ctx.BodyParser(&sf))

		id := ctx.Params("id")
		models.SyncFileUpdateMap(sf, "id=?", id)
		return nil
	})
	r.Get("/file/:id", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		return ctx.JSON(models.SyncFileFindOne("id=?", id))
	})
	r.Get("/files", func(ctx *fiber.Ctx) error {
		pageP := ctx.Query("page")
		page, _ := strconv.Atoi(pageP)
		perPageP := ctx.Query("per_page")
		perPage, _ := strconv.Atoi(perPageP)

		page, perPage = models.Pagination(page, perPage)
		tasks, total, err := models.SyncFileRange(page, perPage, "")
		xerror.Panic(err)

		next, total := models.NextPage(int64(page), int64(perPage), total)
		return ctx.JSON(fiber.Map{
			"total": total,
			"data":  tasks,
			"next":  next,
		})
	})
}
