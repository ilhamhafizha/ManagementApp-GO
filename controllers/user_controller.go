package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ilhamhafizha/ManagementApp-GO/models"
	"github.com/ilhamhafizha/ManagementApp-GO/services"
	"github.com/ilhamhafizha/ManagementApp-GO/utils"
	"github.com/jinzhu/copier"
)

type UserController struct {
	service services.UserService
}

func NewUserController(s services.UserService) *UserController {
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
	var userResponse models.UserResponse
	_ = copier.Copy(&userResponse, user)
	return utils.Success(ctx, "User berhasil didaftarkan", userResponse)
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing data", err.Error())
	}
	user,err := c.service.Login(body.Email, body.Password)
	if err != nil {
		return utils.Unautohorized(ctx, "Login Gagal", err.Error())
	}
	token,_ := utils.GenerateToken(user.InternalID, user.Email, user.Role,user.PublicID)
	refreshToken ,_ := utils.GenerateRefreshToken(user.InternalID)

	var userResponse models.UserResponse
	_ = copier.Copy(&userResponse, user)
	return utils.Success(ctx, "Login Sukses", fiber.Map{
		"access_token": token,
		"refresh_token": refreshToken,
		"user": userResponse,
	})
}

func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := c.service.GetByPublicID(id)
	if err != nil {
		return utils.NotFound(ctx, "User tidak ditemukan", err.Error())
	}
	var userResponse models.UserResponse
	_ = copier.Copy(&userResponse, user)
	if err != nil {
		return utils.BadRequest(ctx, "Gagal parsing data", err.Error())
	}
	return utils.Success(ctx, "User ditemukan", userResponse)
}
