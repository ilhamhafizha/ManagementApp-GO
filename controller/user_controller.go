package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ilhamhafizha/ManagementApp-GO/models"
	"github.com/ilhamhafizha/ManagementApp-GO/services"
	"github.com/ilhamhafizha/ManagementApp-GO/utils"
)

type UserController struct {
	service services.UserService
}

func newUserController(s services.UserService) *UserController {
	return &UserController{service: s}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing data", err.Error())
	}
	if err := c.service.Register(user); err != nil {
		return utils.BadRequest(ctx, "Gagal registrasi user", err.Error())
	}
	return utils.Success(ctx, "User berhasil didaftarkan", nil)
}

