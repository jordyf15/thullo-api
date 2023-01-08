package usecase_test

import (
	"testing"

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

	usecase  card.Usecase
	listRepo *lr.Repository
	cardRepo *cr.Repository
}

var (
	boardID = primitive.NewObjectID()

	list1 = &models.List{
		ID:       primitive.NewObjectID(),
		Title:    "list 1",
		Position: 0,
		BoardID:  boardID,
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

	s.listRepo.On("GetListByID", mock.AnythingOfType("primitive.ObjectID")).Return(list1, nil)
	s.cardRepo.On("GetListCards", mock.AnythingOfType("primitive.ObjectID")).Return([]*models.Card{card1, card2, card3}, nil)
	s.cardRepo.On("Create", mock.AnythingOfType("*models.Card")).Return(nil)

	s.usecase = usecase.NewCardUsecase(s.listRepo, s.cardRepo)
}

func (s *cardUsecaseSuite) TestCreateCardEmptyTitle() {
	err := s.usecase.Create(primitive.NewObjectID(), "")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrCardTitleEmpty.Error(), err.Error())
}

func (s *cardUsecaseSuite) TestCreateCardSuccessful() {
	err := s.usecase.Create(primitive.NewObjectID(), "card 1")

	assert.NoError(s.T(), err)
}
