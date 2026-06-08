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
	AddMember(boardPublicID string, userPublicIDS []string) error
	RemoveMember(boardPublicID string, userPublicIDS []string) error
	GetAllByUserPaginated(userPublicID, filter, sort string, limit, ofset int) ([]models.Board, int64, error)
}

type boardService struct {
	boardRepository repositories.BoardRepository
	userRepository  repositories.UserRepository
	boardMemberRepo repositories.BoardMemberRepository
}

func NewBoardService(
	boardRepository repositories.BoardRepository,
	userRepository repositories.UserRepository,
	boardMemberRepo repositories.BoardMemberRepository,
) BoardService {
	return &boardService{boardRepository, userRepository, boardMemberRepo}
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

func (s *boardService) AddMember(boardPublicID string, userPublicIDS []string) error {
	board, err := s.boardRepository.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}

	var userInternalIDs []uint
	for _, userPublicID := range userPublicIDS {
		user, err := s.userRepository.FindByPublicID(userPublicID)
		if err != nil {
			return errors.New("user not found")
		}
		userInternalIDs = append(userInternalIDs, uint(user.InternalID))
	}

	//cek anggotaan
	existingMember, err := s.boardMemberRepo.GetMembers(string(board.PublicID.String()))
	if err != nil {
		return err
	}

	// cek cepat menggunakan map
	memberMap := make(map[uint]bool)
	for _, member := range existingMember {
		memberMap[uint(member.InternalID)] = true //memberMap[1] = true
	}

	var newMembersIds []uint
	for _, userID := range userInternalIDs {
		if !memberMap[userID] {
			newMembersIds = append(newMembersIds, userID)
		}
	}

	if len(newMembersIds) == 0 {
		return nil
	}

	return s.boardRepository.AddMember(uint(board.InternalID), newMembersIds)
}

func (s *boardService) RemoveMember(boardPublicID string, userPublicIDS []string) error {
	board, err := s.boardRepository.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}

	//validasi user
	var userInternalIDs []uint
	for _, userPublicID := range userPublicIDS {
		user, err := s.userRepository.FindByPublicID(userPublicID)
		if err != nil {
			return errors.New("user not found")
		}
		userInternalIDs = append(userInternalIDs, uint(user.InternalID))
	}

	// cek anggotaan
	existingMember, err := s.boardMemberRepo.GetMembers(string(board.PublicID.String()))
	if err != nil {
		return err
	}

	// cek cepat menggunakan map
	memberMap := make(map[uint]bool)
	for _, member := range existingMember {
		memberMap[uint(member.InternalID)] = true //memberMap[1] = true
	}

	var membersToRemove []uint
	for _, userID := range userInternalIDs {
		if memberMap[userID] {
			membersToRemove = append(membersToRemove, userID)
		}
	}

	return s.boardRepository.RemoveMembers(uint(board.InternalID), membersToRemove)
}

func (s *boardService) GetAllByUserPaginated(userPublicID, filter, sort string, limit, ofset int) ([]models.Board, int64, error) {
	return s.boardRepository.FindAllByUserPaginated(userPublicID, filter, sort, limit, ofset)
}
