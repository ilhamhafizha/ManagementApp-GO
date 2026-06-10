package models

import (
	"github.com/google/uuid"
	"github.com/ilhamhafizha/ManagementApp-GO/models/types"
)

type ListPosition struct {
	InternalID int64     `json:"internal_id" gorm:"column:internal_id;primaryKey;autoIncrement"`
	PublicID   uuid.UUID `json:"public_id" gorm:"column:public_id"`

	BoardID int64 `json:"board_internal_id" gorm:"column:board_internal_id"`

	ListOrder types.UUIDArray `json:"list_order" gorm:"column:list_order"`
}

func (ListPosition) TableName() string {
	return "list_positions"
}
