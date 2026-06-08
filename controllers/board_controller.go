package controllers

import (
	"math"
	"strconv"

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

func (c *BoardController) UpdateBoard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	board := new(models.Board)

	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "Invalid ID Format", err.Error())
	}

	existingBoard, err := c.services.GetByPublicID(publicID)
	if err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Board", err.Error())
	}

	board.InternalID = existingBoard.InternalID
	board.PublicID = existingBoard.PublicID
	board.OwnerPublicID = existingBoard.OwnerPublicID
	board.OwnerID = existingBoard.OwnerID

	if err := c.services.Update(board); err != nil {
		return utils.BadRequest(ctx, "Gagal Update Board", err.Error())
	}
	return utils.Success(ctx, "Board Berhasil Diupdate", board)
}

func (c *BoardController) AddBoardMember(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	var userIDs []string
	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}

	if err := c.services.AddMember(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal Menambahkan Member", err.Error())
	}

	return utils.Success(ctx, "Member Berhasil Ditambahkan", userIDs)
}

func (c *BoardController) RemoveBoardMember(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	var userIDs []string
	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}

	if err := c.services.RemoveMember(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal Menghapus Member", err.Error())
	}

	return utils.Success(ctx, "Member Berhasil dihapus", userIDs)
}

func (c *BoardController) GetMyBoardPaginate(ctx *fiber.Ctx) error {
	jwtUser := ctx.Locals("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	userID := claims["pub_id"].(string)
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page - 1) * limit

	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")

	boards, total, err := c.services.GetAllByUserPaginated(userID, filter, sort, limit, offset)
	if err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	meta := utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		Filter:    filter,
		Sort:      sort,
	}

	return utils.SuccessPagination(ctx, "Data ditemukan", boards, meta)
}
