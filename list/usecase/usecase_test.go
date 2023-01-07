package usecase_test

import (
	"testing"

	br "github.com/jordyf15/thullo-api/board/mocks"
	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/jordyf15/thullo-api/list"
	lr "github.com/jordyf15/thullo-api/list/mocks"
	"github.com/jordyf15/thullo-api/list/usecase"
	"github.com/jordyf15/thullo-api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestListUsecase(t *testing.T) {
	suite.Run(t, new(listUsecaseSuite))
}

var (
	board = &models.Board{
		ID:          primitive.NewObjectID(),
		Title:       "board 1",
		Description: "desc 1",
		OwnerID:     primitive.NewObjectID(),
	}

	list1 = &models.List{
		ID:      primitive.NewObjectID(),
		BoardID: board.ID,
		Title:   "title 1",
	}
)

type listUsecaseSuite struct {
	suite.Suite

	usecase   list.Usecase
	listRepo  *lr.Repository
	boardRepo *br.Repository
}

func (s *listUsecaseSuite) SetupTest() {
	s.boardRepo = new(br.Repository)
	s.listRepo = new(lr.Repository)

	s.boardRepo.On("GetBoardByID", mock.AnythingOfType("primitive.ObjectID")).Return(board, nil)
	s.listRepo.On("GetBoardLists", mock.AnythingOfType("primitive.ObjectID")).Return([]*models.List{list1}, nil)
	s.listRepo.On("Create", mock.AnythingOfType("*models.List")).Return(nil)

	s.usecase = usecase.NewListUsecase(s.listRepo, s.boardRepo)
}

func (s *listUsecaseSuite) TestCreateListEmptyTitle() {
	err := s.usecase.Create(primitive.NewObjectID(), "")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrListTitleEmpty.Error(), err.Error())
}

func (s *listUsecaseSuite) TestCreateListSuccessful() {
	err := s.usecase.Create(primitive.NewObjectID(), "list 1")

	assert.NoError(s.T(), err)
}
