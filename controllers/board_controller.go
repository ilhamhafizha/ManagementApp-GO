package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/ilhamhafizha/ManagementApp-GO/models"
	"github.com/ilhamhafizha/ManagementApp-GO/services"
	"github.com/ilhamhafizha/ManagementApp-GO/utils"
)

type BoardController struct {
	services services.BoardService
}

func NewBoardController(s services.BoardService) *BoardController {
	return &BoardController{services: s}
}

func (c BoardController) CreateBoard(ctx *fiber.Ctx) error {

	var userID uuid.UUID
	var err error

	board := new(models.Board)
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}

	userID, err = uuid.Parse(claims["pub_id"].(string))
	if err != nil {
		return utils.BadRequest(ctx, "Invalid ID Format", err.Error())
	}

	board.OwnerPublicID = userID
	
	if err := c.services.Create(board); err != nil {
		return utils.BadRequest(ctx, "Gagal Membuat Board", err.Error())
	}

	return utils.Success(ctx, "Board Berhasil Dibuat", board)
}
