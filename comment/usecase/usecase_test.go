package usecase_test

import (
	"testing"

	br "github.com/jordyf15/thullo-api/board/mocks"
	bmr "github.com/jordyf15/thullo-api/board_member/mocks"
	cr "github.com/jordyf15/thullo-api/card/mocks"
	"github.com/jordyf15/thullo-api/comment"
	cmr "github.com/jordyf15/thullo-api/comment/mocks"
	"github.com/jordyf15/thullo-api/comment/usecase"
	"github.com/jordyf15/thullo-api/custom_errors"
	lr "github.com/jordyf15/thullo-api/list/mocks"
	"github.com/jordyf15/thullo-api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCommentUsecase(t *testing.T) {
	suite.Run(t, new(commentUsecaseSuite))
}

var (
	board1 = &models.Board{
		ID:         primitive.NewObjectID(),
		Visibility: models.BoardVisibilityPrivate,
	}
	board2 = &models.Board{
		ID:         primitive.NewObjectID(),
		Visibility: models.BoardVisibilityPublic,
	}

	boardMember1 = &models.BoardMember{
		ID:      primitive.NewObjectID(),
		BoardID: board1.ID,
		UserID:  primitive.NewObjectID(),
		Role:    models.MemberRoleAdmin,
	}
	boardMember2 = &models.BoardMember{
		ID:      primitive.NewObjectID(),
		BoardID: board1.ID,
		UserID:  primitive.NewObjectID(),
		Role:    models.MemberRoleMember,
	}

	list1 = &models.List{
		ID:      primitive.NewObjectID(),
		BoardID: board1.ID,
	}
	list2 = &models.List{
		ID:      primitive.NewObjectID(),
		BoardID: board2.ID,
	}
	list3 = &models.List{
		ID:      primitive.NewObjectID(),
		BoardID: primitive.NewObjectID(),
	}

	card1 = &models.Card{
		ID:     primitive.NewObjectID(),
		ListID: list1.ID,
	}
	card2 = &models.Card{
		ID:     primitive.NewObjectID(),
		ListID: list2.ID,
	}
	card3 = &models.Card{
		ID:     primitive.NewObjectID(),
		ListID: primitive.NewObjectID(),
	}

	comment1 = &models.Comment{
		ID:       primitive.NewObjectID(),
		AuthorID: boardMember1.UserID,
		CardID:   card1.ID,
	}
	comment2 = &models.Comment{
		ID:       primitive.NewObjectID(),
		AuthorID: boardMember1.UserID,
		CardID:   card2.ID,
	}
	comment3 = &models.Comment{
		ID:       primitive.NewObjectID(),
		AuthorID: primitive.NewObjectID(),
		CardID:   primitive.NewObjectID(),
	}
)

type commentUsecaseSuite struct {
	suite.Suite

	usecase comment.Usecase

	commentRepo     *cmr.Repository
	boardMemberRepo *bmr.Repository
	boardRepo       *br.Repository
	cardRepo        *cr.Repository
	listRepo        *lr.Repository
}

func (s *commentUsecaseSuite) SetupTest() {
	s.commentRepo = new(cmr.Repository)
	s.boardMemberRepo = new(bmr.Repository)
	s.boardRepo = new(br.Repository)
	s.cardRepo = new(cr.Repository)
	s.listRepo = new(lr.Repository)

	getBoardByID := func(boardID primitive.ObjectID) *models.Board {
		if boardID == board1.ID {
			return board1
		}

		return board2
	}

	getListByID := func(listID primitive.ObjectID) *models.List {
		if listID == list1.ID {
			return list1
		} else if listID == list2.ID {
			return list2
		}

		return list3
	}

	getCardByID := func(cardID primitive.ObjectID) *models.Card {
		if cardID == card1.ID {
			return card1
		} else if cardID == card2.ID {
			return card2
		}

		return card3
	}

	getCommentByID := func(commentID primitive.ObjectID) *models.Comment {
		if commentID == comment1.ID {
			return comment1
		} else if commentID == comment2.ID {
			return comment2
		}

		return comment3
	}

	s.boardRepo.On("GetBoardByID", mock.AnythingOfType("primitive.ObjectID")).Return(getBoardByID, nil)
	s.boardMemberRepo.On("GetBoardMembers", mock.AnythingOfType("primitive.ObjectID")).Return([]*models.BoardMember{boardMember1, boardMember2}, nil)
	s.cardRepo.On("GetCardByID", mock.AnythingOfType("primitive.ObjectID")).Return(getCardByID, nil)
	s.commentRepo.On("Create", mock.AnythingOfType("*models.Comment")).Return(nil)
	s.commentRepo.On("GetCommentByID", mock.AnythingOfType("primitive.ObjectID")).Return(getCommentByID, nil)
	s.commentRepo.On("Update", mock.AnythingOfType("*models.Comment")).Return(nil)
	s.listRepo.On("GetListByID", mock.AnythingOfType("primitive.ObjectID")).Return(getListByID, nil)

	s.usecase = usecase.NewCommentUsecase(s.boardMemberRepo, s.cardRepo, s.commentRepo, s.boardRepo, s.listRepo)
}

func (s *commentUsecaseSuite) TestCreateEmptyComment() {
	err := s.usecase.Create(primitive.NewObjectID(), primitive.NewObjectID(), primitive.NewObjectID(), primitive.NewObjectID(), "")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrCommentEmpty.Error(), err.Error())
	s.commentRepo.AssertNumberOfCalls(s.T(), "Create", 0)
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "GetBoardMembers", 0)
	s.listRepo.AssertNumberOfCalls(s.T(), "GetListByID", 0)
	s.cardRepo.AssertNumberOfCalls(s.T(), "GetCardByID", 0)
}

func (s *commentUsecaseSuite) TestCreateOnPrivateBoardNotAuthorized() {
	err := s.usecase.Create(primitive.NewObjectID(), board1.ID, primitive.NewObjectID(), primitive.NewObjectID(), "comment 1")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrNotAuthorized.Error(), err.Error())
	s.commentRepo.AssertNumberOfCalls(s.T(), "Create", 0)
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "GetBoardMembers", 1)
	s.listRepo.AssertNumberOfCalls(s.T(), "GetListByID", 0)
	s.cardRepo.AssertNumberOfCalls(s.T(), "GetCardByID", 0)
}

func (s *commentUsecaseSuite) TestCreateListNotBelongToBoard() {
	err := s.usecase.Create(boardMember1.UserID, board1.ID, primitive.NewObjectID(), primitive.NewObjectID(), "comment 1")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrRecordNotFound.Error(), err.Error())
	s.commentRepo.AssertNumberOfCalls(s.T(), "Create", 0)
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "GetBoardMembers", 1)
	s.listRepo.AssertNumberOfCalls(s.T(), "GetListByID", 1)
	s.cardRepo.AssertNumberOfCalls(s.T(), "GetCardByID", 0)
}

func (s *commentUsecaseSuite) TestCreateCardNotBelongToList() {
	err := s.usecase.Create(boardMember1.UserID, board1.ID, list1.ID, primitive.NewObjectID(), "comment 1")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrRecordNotFound.Error(), err.Error())
	s.commentRepo.AssertNumberOfCalls(s.T(), "Create", 0)
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "GetBoardMembers", 1)
	s.listRepo.AssertNumberOfCalls(s.T(), "GetListByID", 1)
	s.cardRepo.AssertNumberOfCalls(s.T(), "GetCardByID", 1)
}

func (s *commentUsecaseSuite) TestCreateOnPrivateBoardSuccessful() {
	err := s.usecase.Create(boardMember1.UserID, board1.ID, list1.ID, card1.ID, "comment 1")

	assert.NoError(s.T(), err)
	s.commentRepo.AssertNumberOfCalls(s.T(), "Create", 1)
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "GetBoardMembers", 1)
	s.listRepo.AssertNumberOfCalls(s.T(), "GetListByID", 1)
	s.cardRepo.AssertNumberOfCalls(s.T(), "GetCardByID", 1)
}

func (s *commentUsecaseSuite) TestCreateOnPublicBoardSuccessful() {
	err := s.usecase.Create(primitive.NewObjectID(), board2.ID, list2.ID, card2.ID, "comment 1")

	assert.NoError(s.T(), err)
	s.commentRepo.AssertNumberOfCalls(s.T(), "Create", 1)
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "GetBoardMembers", 0)
	s.listRepo.AssertNumberOfCalls(s.T(), "GetListByID", 1)
	s.cardRepo.AssertNumberOfCalls(s.T(), "GetCardByID", 1)
}

func (s *commentUsecaseSuite) TestUpdateEmptyComment() {
	err := s.usecase.Update(primitive.NewObjectID(), primitive.NewObjectID(), primitive.NewObjectID(), primitive.NewObjectID(), primitive.NewObjectID(), "")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrCommentEmpty.Error(), err.Error())
	s.boardRepo.AssertNumberOfCalls(s.T(), "GetBoardByID", 0)
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "GetBoardMembers", 0)
	s.listRepo.AssertNumberOfCalls(s.T(), "GetListByID", 0)
	s.cardRepo.AssertNumberOfCalls(s.T(), "GetCardByID", 0)
	s.commentRepo.AssertNumberOfCalls(s.T(), "GetCommentByID", 0)
	s.commentRepo.AssertNumberOfCalls(s.T(), "Update", 0)
}

func (s *commentUsecaseSuite) TestUpdateCommentOnPrivateBoardAsNonMember() {
	err := s.usecase.Update(primitive.NewObjectID(), board1.ID, primitive.NewObjectID(), primitive.NewObjectID(), primitive.NewObjectID(), "updated comment 1")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrNotAuthorized.Error(), err.Error())
	s.boardRepo.AssertNumberOfCalls(s.T(), "GetBoardByID", 1)
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "GetBoardMembers", 1)
	s.listRepo.AssertNumberOfCalls(s.T(), "GetListByID", 0)
	s.cardRepo.AssertNumberOfCalls(s.T(), "GetCardByID", 0)
	s.commentRepo.AssertNumberOfCalls(s.T(), "GetCommentByID", 0)
	s.commentRepo.AssertNumberOfCalls(s.T(), "Update", 0)
}

func (s *commentUsecaseSuite) TestUpdateOtherUserComment() {
	err := s.usecase.Update(boardMember2.UserID, board1.ID, list1.ID, card1.ID, comment1.ID, "updated comment 1")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrNotAuthorized.Error(), err.Error())
	s.boardRepo.AssertNumberOfCalls(s.T(), "GetBoardByID", 1)
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "GetBoardMembers", 1)
	s.listRepo.AssertNumberOfCalls(s.T(), "GetListByID", 1)
	s.cardRepo.AssertNumberOfCalls(s.T(), "GetCardByID", 1)
	s.commentRepo.AssertNumberOfCalls(s.T(), "GetCommentByID", 1)
	s.commentRepo.AssertNumberOfCalls(s.T(), "Update", 0)
}

func (s *commentUsecaseSuite) TestUpdateCommentOnListNotInBoard() {
	err := s.usecase.Update(boardMember1.UserID, board1.ID, primitive.NewObjectID(), card1.ID, comment1.ID, "updated comment 1")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrRecordNotFound.Error(), err.Error())
	s.boardRepo.AssertNumberOfCalls(s.T(), "GetBoardByID", 1)
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "GetBoardMembers", 1)
	s.listRepo.AssertNumberOfCalls(s.T(), "GetListByID", 1)
	s.cardRepo.AssertNumberOfCalls(s.T(), "GetCardByID", 0)
	s.commentRepo.AssertNumberOfCalls(s.T(), "GetCommentByID", 0)
	s.commentRepo.AssertNumberOfCalls(s.T(), "Update", 0)
}

func (s *commentUsecaseSuite) TestUpdateCommentOnCardNotInList() {
	err := s.usecase.Update(boardMember1.UserID, board1.ID, list1.ID, primitive.NewObjectID(), comment1.ID, "updated comment 1")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrRecordNotFound.Error(), err.Error())
	s.boardRepo.AssertNumberOfCalls(s.T(), "GetBoardByID", 1)
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "GetBoardMembers", 1)
	s.listRepo.AssertNumberOfCalls(s.T(), "GetListByID", 1)
	s.cardRepo.AssertNumberOfCalls(s.T(), "GetCardByID", 1)
	s.commentRepo.AssertNumberOfCalls(s.T(), "GetCommentByID", 0)
	s.commentRepo.AssertNumberOfCalls(s.T(), "Update", 0)
}

func (s *commentUsecaseSuite) TestUpdateCommentNotInCard() {
	err := s.usecase.Update(boardMember1.UserID, board1.ID, list1.ID, card1.ID, primitive.NewObjectID(), "updated comment 1")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrRecordNotFound.Error(), err.Error())
	s.boardRepo.AssertNumberOfCalls(s.T(), "GetBoardByID", 1)
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "GetBoardMembers", 1)
	s.listRepo.AssertNumberOfCalls(s.T(), "GetListByID", 1)
	s.cardRepo.AssertNumberOfCalls(s.T(), "GetCardByID", 1)
	s.commentRepo.AssertNumberOfCalls(s.T(), "GetCommentByID", 1)
	s.commentRepo.AssertNumberOfCalls(s.T(), "Update", 0)
}

func (s *commentUsecaseSuite) TestUpdateCommentOnPrivateBoardSuccessful() {
	err := s.usecase.Update(boardMember1.UserID, board1.ID, list1.ID, card1.ID, comment1.ID, "updated comment 1")

	assert.NoError(s.T(), err)
	s.boardRepo.AssertNumberOfCalls(s.T(), "GetBoardByID", 1)
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "GetBoardMembers", 1)
	s.listRepo.AssertNumberOfCalls(s.T(), "GetListByID", 1)
	s.cardRepo.AssertNumberOfCalls(s.T(), "GetCardByID", 1)
	s.commentRepo.AssertNumberOfCalls(s.T(), "GetCommentByID", 1)
	s.commentRepo.AssertNumberOfCalls(s.T(), "Update", 1)
}

func (s *commentUsecaseSuite) TestUpdateCommentOnPublicBoardSuccessful() {
	err := s.usecase.Update(boardMember1.UserID, board2.ID, list2.ID, card2.ID, comment2.ID, "updated comment 1")

	assert.NoError(s.T(), err)
	s.boardRepo.AssertNumberOfCalls(s.T(), "GetBoardByID", 1)
	s.listRepo.AssertNumberOfCalls(s.T(), "GetListByID", 1)
	s.cardRepo.AssertNumberOfCalls(s.T(), "GetCardByID", 1)
	s.commentRepo.AssertNumberOfCalls(s.T(), "GetCommentByID", 1)
	s.commentRepo.AssertNumberOfCalls(s.T(), "Update", 1)
}
