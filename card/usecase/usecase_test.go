package usecase_test

import (
	"testing"

	bmr "github.com/jordyf15/thullo-api/board_member/mocks"
	"github.com/jordyf15/thullo-api/card"
	cr "github.com/jordyf15/thullo-api/card/mocks"
	"github.com/jordyf15/thullo-api/card/usecase"
	"github.com/jordyf15/thullo-api/custom_errors"
	lr "github.com/jordyf15/thullo-api/list/mocks"
	"github.com/jordyf15/thullo-api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCardUsecase(t *testing.T) {
	suite.Run(t, new(cardUsecaseSuite))
}

type cardUsecaseSuite struct {
	suite.Suite

	usecase         card.Usecase
	listRepo        *lr.Repository
	cardRepo        *cr.Repository
	boardMemberRepo *bmr.Repository
}

var (
	board1 = &models.Board{
		ID: primitive.NewObjectID(),
	}
	board2 = &models.Board{
		ID: primitive.NewObjectID(),
	}

	boardMember1 = &models.BoardMember{
		ID:      primitive.NewObjectID(),
		UserID:  primitive.NewObjectID(),
		BoardID: board1.ID,
	}
	boardMember2 = &models.BoardMember{
		ID:      primitive.NewObjectID(),
		UserID:  primitive.NewObjectID(),
		BoardID: board2.ID,
	}

	list1 = &models.List{
		ID:       primitive.NewObjectID(),
		Title:    "list 1",
		Position: 0,
		BoardID:  board1.ID,
	}
	list2 = &models.List{
		ID:       primitive.NewObjectID(),
		Title:    "list 2",
		Position: 0,
		BoardID:  primitive.NewObjectID(),
	}

	card1 = &models.Card{
		ID:       primitive.NewObjectID(),
		Title:    "card 1",
		ListID:   list1.ID,
		Position: 0,
	}
	card2 = &models.Card{
		ID:       primitive.NewObjectID(),
		Title:    "card 2",
		ListID:   list1.ID,
		Position: 1,
	}
	card3 = &models.Card{
		ID:       primitive.NewObjectID(),
		Title:    "card 3",
		ListID:   list1.ID,
		Position: 2,
	}
)

func (s *cardUsecaseSuite) SetupTest() {
	s.listRepo = new(lr.Repository)
	s.cardRepo = new(cr.Repository)
	s.boardMemberRepo = new(bmr.Repository)

	getBoardMembers := func(boardID primitive.ObjectID) []*models.BoardMember {
		if boardID == board1.ID {
			return []*models.BoardMember{boardMember1, boardMember2}
		}

		return []*models.BoardMember{}
	}

	getListByID := func(listID primitive.ObjectID) *models.List {
		if listID == list1.ID {
			return list1
		}

		return list2
	}

	s.listRepo.On("GetListByID", mock.AnythingOfType("primitive.ObjectID")).Return(getListByID, nil)
	s.cardRepo.On("GetListCards", mock.AnythingOfType("primitive.ObjectID")).Return([]*models.Card{card1, card2, card3}, nil)
	s.cardRepo.On("Create", mock.AnythingOfType("*models.Card")).Return(nil)
	s.boardMemberRepo.On("GetBoardMembers", mock.AnythingOfType("primitive.ObjectID")).Return(getBoardMembers, nil)

	s.usecase = usecase.NewCardUsecase(s.listRepo, s.cardRepo, s.boardMemberRepo)
}

func (s *cardUsecaseSuite) TestCreateCardEmptyTitle() {
	err := s.usecase.Create(primitive.NewObjectID(), primitive.NewObjectID(), primitive.NewObjectID(), "")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrCardTitleEmpty.Error(), err.Error())
}

func (s *cardUsecaseSuite) TestCreateCardNoBoard() {
	err := s.usecase.Create(primitive.NewObjectID(), board2.ID, primitive.NewObjectID(), "card 1")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrRecordNotFound.Error(), err.Error())
}

func (s *cardUsecaseSuite) TestCreateCardNotAuthorize() {
	err := s.usecase.Create(primitive.NewObjectID(), board1.ID, primitive.NewObjectID(), "card 1")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrNotAuthorized.Error(), err.Error())
}

func (s *cardUsecaseSuite) TestCreateCardListNotBelongToBoard() {
	err := s.usecase.Create(boardMember1.UserID, board1.ID, primitive.NewObjectID(), "card 1")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrRecordNotFound.Error(), err.Error())
}

func (s *cardUsecaseSuite) TestCreateCardSuccessful() {
	err := s.usecase.Create(boardMember1.UserID, board1.ID, list1.ID, "card 1")

	assert.NoError(s.T(), err)
}
