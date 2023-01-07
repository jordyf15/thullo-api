package usecase_test

import (
	"os"
	"sync"
	"testing"

	"github.com/jordyf15/thullo-api/board"
	br "github.com/jordyf15/thullo-api/board/mocks"
	"github.com/jordyf15/thullo-api/board/usecase"
	"github.com/jordyf15/thullo-api/custom_errors"
	sr "github.com/jordyf15/thullo-api/storage/mocks"
	unr "github.com/jordyf15/thullo-api/unsplash/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBoardUsecase(t *testing.T) {
	suite.Run(t, new(boardUsecaseSuite))
}

var (
	img1 *os.File
	img2 *os.File
	img3 *os.File
)

type boardUsecaseSuite struct {
	suite.Suite

	usecase board.Usecase

	boardRepo    *br.Repository
	unsplashRepo *unr.Repository
	storage      *sr.Storage
}

func (s *boardUsecaseSuite) SetupTest() {
	s.boardRepo = new(br.Repository)
	s.unsplashRepo = new(unr.Repository)
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

	s.usecase = usecase.NewBoardUsecase(s.boardRepo, s.unsplashRepo, s.storage)
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
