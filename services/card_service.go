package services

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/ilhamhafizha/ManagementApp-GO/config"
	"github.com/ilhamhafizha/ManagementApp-GO/models"
	"github.com/ilhamhafizha/ManagementApp-GO/models/types"
	"github.com/ilhamhafizha/ManagementApp-GO/repositories"
	"gorm.io/gorm"
)

type CardService interface {
	Create(card *models.Card, listPublicID string) error
	Update(card *models.Card, listPublicID string) error
	Delete(id uint) error

	GetByListID(listPublicID string) ([]models.Card, error)
	GetByID(id uint) (*models.Card, error)
	GetByPublicID(publicID string) (*models.Card, error)
}

type cardService struct {
	cardRepo repositories.CardRepository
	listRepo repositories.ListRepository
	userRepo repositories.UserRepository
}

func NewCardService(cardRepo repositories.CardRepository, listRepo repositories.ListRepository, userRepo repositories.UserRepository) CardService {
	return &cardService{
		cardRepo, listRepo, userRepo,
	}
}

func (s *cardService) Create(card *models.Card, listPublicID string) error {
	// 1. Ambil list dari listPublicID
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found: %w", err)
	}

	// 2. Set list
	card.ListID = list.InternalID

	// 3. Generate public_id jika belum ada
	if card.PublicID == uuid.Nil {
		card.PublicID = uuid.New()
	}
	card.CreatedAt = time.Now()

	// 4. Mulai transaksi
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 5. Simpan card
	if err := tx.Create(card).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create card: %w", err)
	}

	// 6. Update atau buat card_position
	var position models.CardPosition
	if err := tx.Model(&models.CardPosition{}).
		Where("list_internal_id = ?", list.InternalID).
		First(&position).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Buat baru jika belum ada
			position = models.CardPosition{
				PublicID:  uuid.New(),
				ListID:    list.InternalID,
				CardOrder: types.UUIDArray{card.PublicID},
			}

			if err := tx.Create(&position).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create card position: %w", err)
			}
		} else {
			tx.Rollback()
			return fmt.Errorf("failed to get card position: %w", err)
		}
	} else {
		position.CardOrder = append(position.CardOrder, card.PublicID)
		if err := tx.Model(&models.CardPosition{}).
			Where("internal_id = ?", position.InternalID).
			Update("card_order", position.CardOrder).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update card position: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *cardService) Update(card *models.Card, listPublicID string) error {
	// ambil card lama
	existingCard, err := s.cardRepo.FindByPublicID(card.PublicID.String())
	if err != nil {
		return fmt.Errorf("card not found: %w", err)
	}

	//ambil list baru
	newList, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found : %w", err)
	}

	//mulai trx
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// jika pindah list -> hapus dari posisi list lama & tambah ke lict baru

	if existingCard.ListID != newList.InternalID {
		var oldPos models.CardPosition
		if err := tx.Where("list_internal_id = ?", existingCard.ListID).First(&oldPos).Error; err != nil {
			filtered := make(types.UUIDArray, 0, len(oldPos.CardOrder))
			for _, id := range oldPos.CardOrder {
				if id != existingCard.PublicID {
					filtered = append(filtered, id)
				}
			}
			//update

			if err := tx.Model(&models.CardPosition{}).Where("internal_id = ?", oldPos.InternalID).
				Update("card_order", types.UUIDArray(filtered)).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update old card position : %w", err)
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return fmt.Errorf("failed to get old card position : %w", err)
		}

		//tambahkan ke posisi baru
		var newPos models.CardPosition
		res := tx.Where("list_internal_id = ?", newList.InternalID).First(&newPos)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			newPos = models.CardPosition{
				PublicID:  uuid.New(),
				ListID:    newList.InternalID,
				CardOrder: types.UUIDArray{existingCard.PublicID},
			}
			if err := tx.Create(&newPos).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create card position for new list : %w", err)
			}
		} else if res.Error == nil {
			//append
			updateOrder := append(newPos.CardOrder, existingCard.PublicID)
			if err := tx.Model(&models.CardPosition{}).
				Where("internal_id = ?", newPos.InternalID).
				Update("card_order", types.UUIDArray(updateOrder)).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update new card position : %w", err)
			}
		} else {
			tx.Rollback()
			return fmt.Errorf("failed to get new card position : %w", res.Error)
		}
	}

	card.InternalID = existingCard.InternalID
	card.PublicID = existingCard.PublicID
	card.ListID = existingCard.ListID

	if err := tx.Save(card).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update card : %w", err)
	}

	// commit trx
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction : %w", err)
	}

	return nil
}

func (s *cardService) Delete(id uint) error {
	return s.cardRepo.Delete(id)
}

func (s *cardService) GetByListID(listPublicID string) ([]models.Card, error) {

	// ambil list
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return nil, fmt.Errorf("list not found: %w", err)
	}

	// ambil card position
	//position, err := s.cardRepo.FindCardPositionByListID(list.InternalID)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get card position: %w", err)
	//}
	var position *models.CardPosition

	position, err = s.cardRepo.FindCardPositionByListID(list.InternalID)

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to get card position: %w", err)
		}

		position = nil
	}

	// ambil card
	cards, err := s.cardRepo.FindByListID(listPublicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cards: %w", err)
	}

	// sort card berdasarkan position
	if position != nil && len(position.CardOrder) > 0 {
		cards = sortCardByPosition(cards, position.CardOrder)
	}

	return cards, nil
}

func sortCardByPosition(cards []models.Card, order []uuid.UUID) []models.Card {
	orderMap := make(map[uuid.UUID]int)
	for i, id := range order {
		orderMap[id] = i
	}

	defaultIndex := len(order)

	//sorting slice

	sort.SliceStable(cards, func(i, j int) bool {
		idxI, okI := orderMap[cards[i].PublicID]

		if !okI {
			idxI = defaultIndex
		}

		idxJ, okJ := orderMap[cards[j].PublicID]

		if !okJ {
			idxJ = defaultIndex
		}

		if idxI == idxJ {
			return cards[i].CreatedAt.Before(cards[j].CreatedAt)
		}
		return idxI < idxJ
	})

	return cards
}

func (s *cardService) GetByID(id uint) (*models.Card, error) {
	return s.cardRepo.FindByID(id)
}

func (s *cardService) GetByPublicID(publicID string) (*models.Card, error) {
	return s.cardRepo.FindByPublicID(publicID)
}
