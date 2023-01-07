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
		ID:       primitive.NewObjectID(),
		BoardID:  board.ID,
		Title:    "title 1",
		Position: 0,
	}
	list2 = &models.List{
		ID:       primitive.NewObjectID(),
		BoardID:  board.ID,
		Title:    "title 2",
		Position: 1,
	}
	list3 = &models.List{
		ID:       primitive.NewObjectID(),
		BoardID:  board.ID,
		Title:    "title 3",
		Position: 2,
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

	// need to reset list position
	list1.Position = 0
	list2.Position = 1
	list3.Position = 2

	getListByID := func(listID primitive.ObjectID) *models.List {
		if listID == list3.ID {
			return list3
		}

		return list1
	}

	s.boardRepo.On("GetBoardByID", mock.AnythingOfType("primitive.ObjectID")).Return(board, nil)
	s.listRepo.On("GetBoardLists", mock.AnythingOfType("primitive.ObjectID")).Return([]*models.List{list1, list2, list3}, nil)
	s.listRepo.On("Create", mock.AnythingOfType("*models.List")).Return(nil)
	s.listRepo.On("GetListByID", mock.AnythingOfType("primitive.ObjectID")).Return(getListByID, nil)
	s.listRepo.On("UpdateList", mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("*models.List")).Return(nil)

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

func (s *listUsecaseSuite) TestUpdateTitleEmptyTitle() {
	err := s.usecase.UpdateTitle(primitive.NewObjectID(), "")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrListTitleEmpty.Error(), err.Error())
}

func (s *listUsecaseSuite) TestUpdateTitleSuccessful() {
	err := s.usecase.UpdateTitle(primitive.NewObjectID(), "new title")

	assert.NoError(s.T(), err)
}

func (s *listUsecaseSuite) TestUpdatePositionUpward() {
	err := s.usecase.UpdatePosition(primitive.NewObjectID(), list1.ID, 1)

	assert.NoError(s.T(), err)
	s.listRepo.AssertNumberOfCalls(s.T(), "UpdateList", 2)
	assert.Equal(s.T(), 1, list1.Position)
	assert.Equal(s.T(), 0, list2.Position)
	assert.Equal(s.T(), 2, list3.Position)
}

func (s *listUsecaseSuite) TestUpdatePositionDownward() {
	err := s.usecase.UpdatePosition(primitive.NewObjectID(), list3.ID, 1)

	assert.NoError(s.T(), err)
	s.listRepo.AssertNumberOfCalls(s.T(), "UpdateList", 2)
	assert.Equal(s.T(), 0, list1.Position)
	assert.Equal(s.T(), 1, list3.Position)
	assert.Equal(s.T(), 2, list2.Position)
}
