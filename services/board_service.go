package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ilhamhafizha/ManagementApp-GO/models"
	"github.com/ilhamhafizha/ManagementApp-GO/repositories"
)

type BoardService interface {
	Create(board *models.Board) error
	Update(board *models.Board) error
	GetByPublicID(publicID string) (*models.Board, error)
}

type boardService struct {
	boardRepository repositories.BoardRepository
	userRepository  repositories.UserRepository
}

func NewBoardService(
	boardRepository repositories.BoardRepository,
	userRepository repositories.UserRepository,
) BoardService {
	return &boardService{boardRepository, userRepository}
}

func (s *boardService) Create(board *models.Board) error {
	user, err := s.userRepository.FindByPublicID(board.OwnerPublicID.String())
	if err != nil {
		return errors.New("user not found")
	}
	board.PublicID = uuid.New()
	board.OwnerID = user.InternalID
	return s.boardRepository.Create(board)
}

func (s *boardService) Update(board *models.Board) error {
	return s.boardRepository.Update(board)
}

func (s *boardService) GetByPublicID(publicID string) (*models.Board, error) {
	return s.boardRepository.FindByPublicID(publicID)
}
