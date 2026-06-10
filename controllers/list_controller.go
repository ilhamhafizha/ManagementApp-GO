package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ilhamhafizha/ManagementApp-GO/models"
	"github.com/ilhamhafizha/ManagementApp-GO/services"
	"github.com/ilhamhafizha/ManagementApp-GO/utils"
)

type ListController struct {
	service services.ListService
}

func NewListController(s services.ListService) *ListController {
	return &ListController{service: s}
}

func (c *ListController) CreateList(ctx *fiber.Ctx) error {
	list := new(models.List)
	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Gagal Membaca Request", err.Error())
	}
	if err := c.service.Create(list); err != nil {
		return utils.BadRequest(ctx, "Gagal Membuat List", err.Error())
	}

	return utils.Success(ctx, "List Berhasil Dibuat", list)
}

func (c *ListController) UpdateList(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	list := new(models.List)
	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Gagal Membaca Request", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	existingList, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan", err.Error())
	}

	list.InternalID = existingList.InternalID
	list.PublicID = existingList.PublicID

	if err := c.service.Update(list); err != nil {
		return utils.BadRequest(ctx, "Gagal Memperbaharui List", err.Error())
	}

	updatedList, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan", err.Error())
	}

	return utils.Success(ctx, "List Memperbaharui Dibuat", updatedList)
}

func (c *ListController) GetListOnBoard(ctx *fiber.Ctx) error {
	boardPublicID := ctx.Params("board_id")
	if _, err := uuid.Parse(boardPublicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	lists, err := c.service.GetByBoardID(boardPublicID)
	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan", err.Error())
	}

	return utils.Success(ctx, "List Ditemukan", lists)
}

func (c *ListController) DeleteList(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	list, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan", err.Error())
	}

	if err := c.service.Delete(uint(list.InternalID)); err != nil {
		return utils.InternalServerError(ctx, "Gagal Menghapus List", err.Error())
	}

	return utils.Success(ctx, "List Berhasil Dihapus", list)
}

func (c *ListController) UpdateListPosition(ctx *fiber.Ctx) error {
	boardID := ctx.Params("board_id")
	if _, err := uuid.Parse(boardID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	var positionUUID []uuid.UUID
	if err := ctx.BodyParser(&positionUUID); err != nil {
		//jika gagal, coba parse sebagai array of string
		var positionString []string
		if err := ctx.BodyParser(&positionString); err != nil {
			return utils.BadRequest(ctx, "Invalid position format", err.Error())
		}
		// konversi string ke uuid
		for _, s := range positionString {
			u, err := uuid.Parse(s)
			if err != nil {
				return utils.BadRequest(ctx, "Invalid position format", err.Error())
			}
			positionUUID = append(positionUUID, u)
		}
	}
	if err := c.service.UpdatePositions(boardID, positionUUID); err != nil {
		return utils.InternalServerError(ctx, "Gagal Memperbaharui Posisi List", err.Error())
	}
	return utils.Success(ctx, "Posisi List Berhasil Diperbaharui", nil)
}
