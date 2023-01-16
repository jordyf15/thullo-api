package usecase_test

import (
	"testing"

	br "github.com/jordyf15/thullo-api/board/mocks"
	bmr "github.com/jordyf15/thullo-api/board_member/mocks"
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
	board1 = &models.Board{
		ID:          primitive.NewObjectID(),
		Title:       "board 1",
		Description: "desc 1",
		OwnerID:     primitive.NewObjectID(),
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
		BoardID: board1.ID,
	}

	list1 = &models.List{
		ID:       primitive.NewObjectID(),
		BoardID:  board1.ID,
		Title:    "title 1",
		Position: 0,
	}
	list2 = &models.List{
		ID:       primitive.NewObjectID(),
		BoardID:  board1.ID,
		Title:    "title 2",
		Position: 1,
	}
	list3 = &models.List{
		ID:       primitive.NewObjectID(),
		BoardID:  board1.ID,
		Title:    "title 3",
		Position: 2,
	}
	list4 = &models.List{
		ID:       primitive.NewObjectID(),
		BoardID:  board2.ID,
		Title:    "title 1",
		Position: 0,
	}
)

type listUsecaseSuite struct {
	suite.Suite

	usecase         list.Usecase
	listRepo        *lr.Repository
	boardRepo       *br.Repository
	boardMemberRepo *bmr.Repository
}

func (s *listUsecaseSuite) SetupTest() {
	s.boardRepo = new(br.Repository)
	s.listRepo = new(lr.Repository)
	s.boardMemberRepo = new(bmr.Repository)

	// need to reset list position
	list1.Position = 0
	list2.Position = 1
	list3.Position = 2

	getListByID := func(listID primitive.ObjectID) *models.List {
		if listID == list3.ID {
			return list3
		} else if listID == list1.ID {
			return list1
		} else if listID == list2.ID {
			return list2
		}

		return list4
	}

	getBoardMembers := func(boardID primitive.ObjectID) []*models.BoardMember {
		if boardID == board1.ID {
			return []*models.BoardMember{boardMember1, boardMember2}
		}

		return []*models.BoardMember{}
	}

	s.listRepo.On("GetBoardLists", mock.AnythingOfType("primitive.ObjectID")).Return([]*models.List{list1, list2, list3}, nil)
	s.listRepo.On("Create", mock.AnythingOfType("*models.List")).Return(nil)
	s.listRepo.On("GetListByID", mock.AnythingOfType("primitive.ObjectID")).Return(getListByID, nil)
	s.listRepo.On("UpdateList", mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("*models.List")).Return(nil)
	s.boardMemberRepo.On("GetBoardMembers", mock.AnythingOfType("primitive.ObjectID")).Return(getBoardMembers, nil)

	s.usecase = usecase.NewListUsecase(s.listRepo, s.boardRepo, s.boardMemberRepo)
}

func (s *listUsecaseSuite) TestCreateListEmptyTitle() {
	err := s.usecase.Create(primitive.NewObjectID(), primitive.NewObjectID(), "")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrListTitleEmpty.Error(), err.Error())
}

func (s *listUsecaseSuite) TestCreateListBoardDoesNotExist() {
	err := s.usecase.Create(primitive.NewObjectID(), primitive.NewObjectID(), "todo 1")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrRecordNotFound.Error(), err.Error())
}

func (s *listUsecaseSuite) TestCreateListUserNotAuthorized() {
	err := s.usecase.Create(primitive.NewObjectID(), board1.ID, "todo 1")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrNotAuthorized.Error(), err.Error())
}

func (s *listUsecaseSuite) TestCreateListSuccessful() {
	err := s.usecase.Create(boardMember1.UserID, board1.ID, "todo 1")

	assert.NoError(s.T(), err)
}

func (s *listUsecaseSuite) TestUpdateTitleEmptyTitle() {
	err := s.usecase.UpdateTitle(primitive.NewObjectID(), primitive.NewObjectID(), primitive.NewObjectID(), "")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrListTitleEmpty.Error(), err.Error())
	s.listRepo.AssertNumberOfCalls(s.T(), "UpdateList", 0)
}

func (s *listUsecaseSuite) TestUpdateTitleBoardNotFound() {
	err := s.usecase.UpdateTitle(primitive.NewObjectID(), primitive.NewObjectID(), primitive.NewObjectID(), "todo 1 updated")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrRecordNotFound.Error(), err.Error())
	s.listRepo.AssertNumberOfCalls(s.T(), "UpdateList", 0)
}

func (s *listUsecaseSuite) TestUpdateTitleUserNotAuthorized() {
	err := s.usecase.UpdateTitle(primitive.NewObjectID(), board1.ID, primitive.NewObjectID(), "todo 1 updated")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrNotAuthorized.Error(), err.Error())
	s.listRepo.AssertNumberOfCalls(s.T(), "UpdateList", 0)
}

func (s *listUsecaseSuite) TestUpdateTitleListNotBelongToBoard() {
	err := s.usecase.UpdateTitle(boardMember1.UserID, board1.ID, primitive.NewObjectID(), "todo1 updated")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrRecordNotFound.Error(), err.Error())
	s.listRepo.AssertNumberOfCalls(s.T(), "UpdateList", 0)
}

func (s *listUsecaseSuite) TestUpdateTitleSuccessful() {
	err := s.usecase.UpdateTitle(boardMember1.UserID, board1.ID, list1.ID, "todo 1 updated")

	assert.NoError(s.T(), err)
	s.listRepo.AssertNumberOfCalls(s.T(), "UpdateList", 1)
}

func (s *listUsecaseSuite) TestUpdatePositionTooLow() {
	err := s.usecase.UpdatePosition(boardMember1.UserID, board1.ID, list1.ID, -1)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrListPositionTooLow.Error(), err.Error())
	s.listRepo.AssertNumberOfCalls(s.T(), "UpdateList", 0)
}

func (s *listUsecaseSuite) TestUpdatePositionTooHigh() {
	err := s.usecase.UpdatePosition(boardMember1.UserID, board1.ID, list1.ID, 4)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrListPositionTooHigh.Error(), err.Error())
	s.listRepo.AssertNumberOfCalls(s.T(), "UpdateList", 0)
}

func (s *listUsecaseSuite) TestUpdatePositionBoardNotFound() {
	err := s.usecase.UpdatePosition(boardMember1.UserID, board2.ID, list1.ID, 2)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrRecordNotFound.Error(), err.Error())
	s.listRepo.AssertNumberOfCalls(s.T(), "UpdateList", 0)
}

func (s *listUsecaseSuite) TestUpdatePositionNotAuthorized() {
	err := s.usecase.UpdatePosition(primitive.NewObjectID(), board1.ID, list1.ID, 2)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrNotAuthorized.Error(), err.Error())
	s.listRepo.AssertNumberOfCalls(s.T(), "UpdateList", 0)
}

func (s *listUsecaseSuite) TestUpdateListNotBelongToBoard() {
	err := s.usecase.UpdatePosition(boardMember1.UserID, board1.ID, primitive.NewObjectID(), 2)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrRecordNotFound.Error(), err.Error())
	s.listRepo.AssertNumberOfCalls(s.T(), "UpdateList", 0)
}

func (s *listUsecaseSuite) TestUpdatePositionUpward() {
	err := s.usecase.UpdatePosition(boardMember1.UserID, board1.ID, list1.ID, 1)

	assert.NoError(s.T(), err)
	s.listRepo.AssertNumberOfCalls(s.T(), "UpdateList", 2)
	assert.Equal(s.T(), 1, list1.Position)
	assert.Equal(s.T(), 0, list2.Position)
	assert.Equal(s.T(), 2, list3.Position)
}

func (s *listUsecaseSuite) TestUpdatePositionDownward() {
	err := s.usecase.UpdatePosition(boardMember1.UserID, board1.ID, list3.ID, 1)

	assert.NoError(s.T(), err)
	s.listRepo.AssertNumberOfCalls(s.T(), "UpdateList", 2)
	assert.Equal(s.T(), 0, list1.Position)
	assert.Equal(s.T(), 1, list3.Position)
	assert.Equal(s.T(), 2, list2.Position)
}
