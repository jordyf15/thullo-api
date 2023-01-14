package usecase_test

import (
	"os"
	"sync"
	"testing"

	"github.com/jordyf15/thullo-api/board"
	br "github.com/jordyf15/thullo-api/board/mocks"
	"github.com/jordyf15/thullo-api/board/usecase"
	bmr "github.com/jordyf15/thullo-api/board_member/mocks"
	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/jordyf15/thullo-api/models"
	sr "github.com/jordyf15/thullo-api/storage/mocks"
	unr "github.com/jordyf15/thullo-api/unsplash/mocks"
	ur "github.com/jordyf15/thullo-api/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBoardUsecase(t *testing.T) {
	suite.Run(t, new(boardUsecaseSuite))
}

var (
	img1         *os.File
	img2         *os.File
	img3         *os.File
	requesterID1 = primitive.NewObjectID()
	requesterID2 = primitive.NewObjectID()
	newMemberID1 = primitive.NewObjectID()
	newMemberID2 = primitive.NewObjectID()

	board1 = &models.Board{
		ID: primitive.NewObjectID(),
	}
	boardMember1 = &models.BoardMember{
		ID:      primitive.NewObjectID(),
		UserID:  requesterID1,
		BoardID: board1.ID,
		Role:    models.MemberRoleAdmin,
	}
	boardMember2 = &models.BoardMember{
		ID:      primitive.NewObjectID(),
		UserID:  newMemberID1,
		BoardID: board1.ID,
		Role:    models.MemberRoleMember,
	}
	user1 = &models.User{
		ID: newMemberID1,
	}
)

type boardUsecaseSuite struct {
	suite.Suite

	usecase board.Usecase

	boardRepo       *br.Repository
	unsplashRepo    *unr.Repository
	boardMemberRepo *bmr.Repository
	userRepo        *ur.Repository
	storage         *sr.Storage
}

func (s *boardUsecaseSuite) SetupTest() {
	s.boardRepo = new(br.Repository)
	s.unsplashRepo = new(unr.Repository)
	s.userRepo = new(ur.Repository)
	s.boardMemberRepo = new(bmr.Repository)
	s.storage = new(sr.Storage)

	img1, _ = os.Create("image1.jpg")
	img2, _ = os.Create("image2.jpg")
	img3, _ = os.Create("image3.jpg")

	getBoardMembers := func(boardID primitive.ObjectID) []*models.BoardMember {
		if boardID == board1.ID {
			return []*models.BoardMember{boardMember1, boardMember2}
		}

		return []*models.BoardMember{}
	}

	s.unsplashRepo.On("GetImagesForID", mock.AnythingOfType("string"), mock.AnythingOfType("float64")).Return([]*os.File{img1, img2, img3}, nil)
	s.storage.On("UploadFile", mock.AnythingOfType("chan<- error"), mock.AnythingOfType("*sync.WaitGroup"), mock.AnythingOfType("*models.Image"), mock.AnythingOfType("*os.File"), mock.AnythingOfType("map[string]string")).Run(func(args mock.Arguments) {
		arg1 := args[0].(chan<- error)
		arg1 <- nil
		arg2 := args[1].(*sync.WaitGroup)
		arg2.Done()
	})
	s.boardRepo.On("Create", mock.AnythingOfType("*models.Board")).Return(nil)
	s.boardRepo.On("GetBoardByID", mock.AnythingOfType("primitive.ObjectID")).Return(board1, nil)
	s.boardRepo.On("Update", mock.AnythingOfType("*models.Board")).Return(nil)
	s.userRepo.On("GetByID", mock.AnythingOfType("primitive.ObjectID")).Return(user1, nil)
	s.boardMemberRepo.On("Create", mock.AnythingOfType("*models.BoardMember")).Return(nil)
	s.boardMemberRepo.On("GetBoardMembers", mock.AnythingOfType("primitive.ObjectID")).Return(getBoardMembers, nil)
	s.boardMemberRepo.On("UpdateBoardMemberRole", mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("models.MemberRole")).Return(nil)
	s.boardMemberRepo.On("DeleteBoardMemberByID", mock.AnythingOfType("primitive.ObjectID")).Return(nil)

	s.usecase = usecase.NewBoardUsecase(s.boardRepo, s.unsplashRepo, s.boardMemberRepo, s.userRepo, s.storage)
}

func (s *boardUsecaseSuite) AfterTest(suiteName, testName string) {
	img1.Close()
	img2.Close()
	img3.Close()

	os.Remove("image1.jpg")
	os.Remove("image2.jpg")
	os.Remove("image3.jpg")
}

func (s *boardUsecaseSuite) TestCreateInvalidBoardData() {
	err := s.usecase.Create(primitive.NewObjectID(), "", "secret", map[string]interface{}{
		"source":   "imgur",
		"fp_y":     float64(100),
		"photo_id": "picture-1",
	})

	expectedErrors := &custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrInvalidCoverSource, custom_errors.ErrUnsplashFocalPointYTooHigh, custom_errors.ErrBoardTitleEmpty, custom_errors.ErrBoardInvalidVisibility}}

	assert.Equal(s.T(), expectedErrors.Error(), err.Error())
	s.boardRepo.AssertNumberOfCalls(s.T(), "Create", 0)
}

func (s *boardUsecaseSuite) TestCreateSuccessful() {
	err := s.usecase.Create(primitive.NewObjectID(), "Board 1", "public", map[string]interface{}{
		"source":   "unsplash",
		"fp_y":     float64(0.5),
		"photo_id": "picture-1",
	})

	assert.NoError(s.T(), err)
	s.boardRepo.AssertNumberOfCalls(s.T(), "Create", 1)
}

func (s *boardUsecaseSuite) TestUpdateBoardVisibilityInvalidVisibility() {
	err := s.usecase.UpdateVisibility(primitive.NewObjectID(), primitive.NewObjectID(), "visible")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrBoardInvalidVisibility.Error(), err.Error())
	s.boardRepo.AssertNumberOfCalls(s.T(), "Update", 0)
}

func (s *boardUsecaseSuite) TestUpdateBoardVisibilityAsNonMember() {
	err := s.usecase.UpdateVisibility(primitive.NewObjectID(), board1.ID, "public")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrNotAuthorized.Error(), err.Error())
	s.boardRepo.AssertNumberOfCalls(s.T(), "Update", 0)
}

func (s *boardUsecaseSuite) TestUpdateBoardVisibilityAsMember() {
	err := s.usecase.UpdateVisibility(boardMember2.UserID, board1.ID, "public")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrNotAuthorized.Error(), err.Error())
	s.boardRepo.AssertNumberOfCalls(s.T(), "Update", 0)
}

func (s *boardUsecaseSuite) TestUpdateBoardVisibilitySuccessful() {
	err := s.usecase.UpdateVisibility(boardMember1.UserID, board1.ID, "public")

	assert.NoError(s.T(), err)
	s.boardRepo.AssertNumberOfCalls(s.T(), "Update", 1)
}

func (s *boardUsecaseSuite) TestUpdateBoardTitleEmptyTitle() {
	err := s.usecase.UpdateTitle(primitive.NewObjectID(), primitive.NewObjectID(), "")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrBoardTitleEmpty.Error(), err.Error())
	s.boardRepo.AssertNumberOfCalls(s.T(), "Update", 0)
}

func (s *boardUsecaseSuite) TestUpdateBoardTitleAsNonMember() {
	err := s.usecase.UpdateTitle(primitive.NewObjectID(), board1.ID, "updated title")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrNotAuthorized.Error(), err.Error())
	s.boardRepo.AssertNumberOfCalls(s.T(), "Update", 0)
}

func (s *boardUsecaseSuite) TestUpdateBoardTitleAsMember() {
	err := s.usecase.UpdateTitle(boardMember2.UserID, board1.ID, "updated title")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrNotAuthorized.Error(), err.Error())
	s.boardRepo.AssertNumberOfCalls(s.T(), "Update", 0)
}

func (s *boardUsecaseSuite) TestUpdateBoardTitleSuccessful() {
	err := s.usecase.UpdateTitle(boardMember1.UserID, board1.ID, "updated title")

	assert.NoError(s.T(), err)
	s.boardRepo.AssertNumberOfCalls(s.T(), "Update", 1)
}

func (s *boardUsecaseSuite) TestUpdateBoardDescriptionAsNonMember() {
	err := s.usecase.UpdateDescription(primitive.NewObjectID(), board1.ID, "updated description")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrNotAuthorized.Error(), err.Error())
	s.boardRepo.AssertNumberOfCalls(s.T(), "Update", 0)
}

func (s *boardUsecaseSuite) TestUpdateBoardDescriptionSuccessful() {
	err := s.usecase.UpdateDescription(boardMember2.UserID, board1.ID, "updated description")

	assert.NoError(s.T(), err)
	s.boardRepo.AssertNumberOfCalls(s.T(), "Update", 1)
}

func (s *boardUsecaseSuite) TestAddMemberNotAuthorized() {
	err := s.usecase.AddMember(requesterID2, board1.ID, newMemberID2)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrNotAuthorized.Error(), err.Error())
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "Create", 0)
}

func (s *boardUsecaseSuite) TestAddedMemberIsAlreadyMember() {
	err := s.usecase.AddMember(requesterID1, board1.ID, newMemberID1)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrUserIsAlreadyBoardMember.Error(), err.Error())
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "Create", 0)
}

func (s *boardUsecaseSuite) TestAddMemberSuccessful() {
	err := s.usecase.AddMember(requesterID1, board1.ID, newMemberID2)

	assert.NoError(s.T(), err)
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "Create", 1)
}

func (s *boardUsecaseSuite) TestUpdateMemberRoleInvalidRole() {
	err := s.usecase.UpdateMemberRole(primitive.NewObjectID(), primitive.NewObjectID(), primitive.NewObjectID(), "master")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrInvalidBoardMemberRole.Error(), err.Error())
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "UpdateBoardMemberRole", 0)
}

func (s *boardUsecaseSuite) TestUpdateMemberRoleNoMembers() {
	err := s.usecase.UpdateMemberRole(primitive.NewObjectID(), primitive.NewObjectID(), primitive.NewObjectID(), "admin")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrRecordNotFound.Error(), err.Error())
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "UpdateBoardMemberRole", 0)
}

func (s *boardUsecaseSuite) TestUpdateMemberRoleAsNonMember() {
	err := s.usecase.UpdateMemberRole(primitive.NewObjectID(), board1.ID, primitive.NewObjectID(), "member")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrNotAuthorized.Error(), err.Error())
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "UpdateBoardMemberRole", 0)
}

func (s *boardUsecaseSuite) TestUpdateMemberRoleAsMember() {
	err := s.usecase.UpdateMemberRole(boardMember2.UserID, board1.ID, primitive.NewObjectID(), "member")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrNotAuthorized.Error(), err.Error())
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "UpdateBoardMemberRole", 0)
}

func (s *boardUsecaseSuite) TestUpdateMemberRoleForNonMember() {
	err := s.usecase.UpdateMemberRole(boardMember1.UserID, board1.ID, primitive.NewObjectID(), "member")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrRecordNotFound.Error(), err.Error())
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "UpdateBoardMemberRole", 0)
}

func (s *boardUsecaseSuite) TestUpdateMemberRoleDemoteLastAdmin() {
	err := s.usecase.UpdateMemberRole(boardMember1.UserID, board1.ID, boardMember1.UserID, "member")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrBoardMustHaveAnAdmin.Error(), err.Error())
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "UpdateBoardMemberRole", 0)
}

func (s *boardUsecaseSuite) TestUpdateMemberRoleSuccessful() {
	err := s.usecase.UpdateMemberRole(boardMember1.UserID, board1.ID, boardMember2.UserID, "admin")

	assert.NoError(s.T(), err)
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "UpdateBoardMemberRole", 1)
}

func (s *boardUsecaseSuite) TestDeleteMemberNoMembers() {
	err := s.usecase.DeleteMember(primitive.NewObjectID(), primitive.NewObjectID(), primitive.NewObjectID())

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrRecordNotFound.Error(), err.Error())
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "DeleteBoardMemberByID", 0)
}

func (s *boardUsecaseSuite) TestDeleteMemberAsNonMember() {
	err := s.usecase.DeleteMember(primitive.NewObjectID(), board1.ID, primitive.NewObjectID())

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrNotAuthorized.Error(), err.Error())
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "DeleteBoardMemberByID", 0)
}

func (s *boardUsecaseSuite) TestDeleteMemberAsMember() {
	err := s.usecase.DeleteMember(boardMember2.UserID, board1.ID, primitive.NewObjectID())

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrNotAuthorized.Error(), err.Error())
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "DeleteBoardMemberByID", 0)
}

func (s *boardUsecaseSuite) TestDeleteMemberLastAdmin() {
	err := s.usecase.DeleteMember(boardMember1.UserID, board1.ID, boardMember1.UserID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrBoardMustHaveAnAdmin.Error(), err.Error())
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "DeleteBoardMemberByID", 0)
}

func (s *boardUsecaseSuite) TestDeleteMemberSuccessful() {
	err := s.usecase.DeleteMember(boardMember1.UserID, board1.ID, boardMember2.UserID)

	assert.NoError(s.T(), err)
	s.boardMemberRepo.AssertNumberOfCalls(s.T(), "DeleteBoardMemberByID", 1)
}
