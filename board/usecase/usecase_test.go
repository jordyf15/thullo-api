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
	}
	boardMember2 = &models.BoardMember{
		ID:      primitive.NewObjectID(),
		UserID:  newMemberID1,
		BoardID: board1.ID,
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

	s.unsplashRepo.On("GetImagesForID", mock.AnythingOfType("string"), mock.AnythingOfType("float64")).Return([]*os.File{img1, img2, img3}, nil)
	s.storage.On("UploadFile", mock.AnythingOfType("chan<- error"), mock.AnythingOfType("*sync.WaitGroup"), mock.AnythingOfType("*models.Image"), mock.AnythingOfType("*os.File"), mock.AnythingOfType("map[string]string")).Run(func(args mock.Arguments) {
		arg1 := args[0].(chan<- error)
		arg1 <- nil
		arg2 := args[1].(*sync.WaitGroup)
		arg2.Done()
	})
	s.boardRepo.On("Create", mock.AnythingOfType("*models.Board")).Return(nil)
	s.boardRepo.On("GetBoardByID", mock.AnythingOfType("primitive.ObjectID")).Return(board1, nil)
	s.userRepo.On("GetByID", mock.AnythingOfType("primitive.ObjectID")).Return(user1, nil)
	s.boardMemberRepo.On("Create", mock.AnythingOfType("*models.BoardMember")).Return(nil)
	s.boardMemberRepo.On("GetBoardMembers", mock.AnythingOfType("primitive.ObjectID")).Return([]*models.BoardMember{boardMember1, boardMember2}, nil)

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

func (s *boardUsecaseSuite) TestAddMemberNotAuthorized() {
	err := s.usecase.AddMember(requesterID2, board1.ID, newMemberID2)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrNotAuthorized.Error(), err.Error())
}

func (s *boardUsecaseSuite) TestAddedMemberIsAlreadyMember() {
	err := s.usecase.AddMember(requesterID1, board1.ID, newMemberID1)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrUserIsAlreadyBoardMember.Error(), err.Error())
}

func (s *boardUsecaseSuite) TestAddMemberSuccessful() {
	err := s.usecase.AddMember(requesterID1, board1.ID, newMemberID2)

	assert.NoError(s.T(), err)
}
