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

	card1 = &models.Card{
		ID: primitive.NewObjectID(),
	}
)

type commentUsecaseSuite struct {
	suite.Suite

	usecase comment.Usecase

	commentRepo     *cmr.Repository
	boardMemberRepo *bmr.Repository
	boardRepo       *br.Repository
	cardRepo        *cr.Repository
}

func (s *commentUsecaseSuite) SetupTest() {
	s.commentRepo = new(cmr.Repository)
	s.boardMemberRepo = new(bmr.Repository)
	s.boardRepo = new(br.Repository)
	s.cardRepo = new(cr.Repository)

	getBoardByID := func(boardID primitive.ObjectID) *models.Board {
		if boardID == board1.ID {
			return board1
		}

		return board2
	}

	s.boardRepo.On("GetBoardByID", mock.AnythingOfType("primitive.ObjectID")).Return(getBoardByID, nil)
	s.boardMemberRepo.On("GetBoardMembers", mock.AnythingOfType("primitive.ObjectID")).Return([]*models.BoardMember{boardMember1, boardMember2}, nil)
	s.cardRepo.On("GetCardByID", mock.AnythingOfType("primitive.ObjectID")).Return(card1, nil)
	s.commentRepo.On("Create", mock.AnythingOfType("*models.Comment")).Return(nil)

	s.usecase = usecase.NewCommentUsecase(s.boardMemberRepo, s.cardRepo, s.commentRepo, s.boardRepo)
}

func (s *commentUsecaseSuite) TestCreateEmptyComment() {
	err := s.usecase.Create(primitive.NewObjectID(), primitive.NewObjectID(), primitive.NewObjectID(), "")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrCommentEmpty.Error(), err.Error())
	s.commentRepo.AssertNumberOfCalls(s.T(), "Create", 0)
}

func (s *commentUsecaseSuite) TestCreateOnPrivateBoardNotAuthorized() {
	err := s.usecase.Create(primitive.NewObjectID(), board1.ID, primitive.NewObjectID(), "comment 1")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrNotAuthorized.Error(), err.Error())
	s.commentRepo.AssertNumberOfCalls(s.T(), "Create", 0)
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "GetBoardMembers", 1)
}

func (s *commentUsecaseSuite) TestCreateOnPrivateBoardSuccessful() {
	err := s.usecase.Create(boardMember1.UserID, board1.ID, primitive.NewObjectID(), "comment 1")

	assert.NoError(s.T(), err)
	s.commentRepo.AssertNumberOfCalls(s.T(), "Create", 1)
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "GetBoardMembers", 1)
}

func (s *commentUsecaseSuite) TestCreateOnPublicBoardSuccessful() {
	err := s.usecase.Create(primitive.NewObjectID(), board2.ID, primitive.NewObjectID(), "comment 1")

	assert.NoError(s.T(), err)
	s.commentRepo.AssertNumberOfCalls(s.T(), "Create", 1)
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "GetBoardMembers", 0)
}
