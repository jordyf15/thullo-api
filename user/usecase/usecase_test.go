package usecase_test

import (
	"sync"
	"testing"

	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/jordyf15/thullo-api/models"
	or "github.com/jordyf15/thullo-api/oauth/mocks"
	sr "github.com/jordyf15/thullo-api/storage/mocks"
	tr "github.com/jordyf15/thullo-api/token/mocks"
	"github.com/jordyf15/thullo-api/user"
	ur "github.com/jordyf15/thullo-api/user/mocks"
	"github.com/jordyf15/thullo-api/user/usecase"
	ubr "github.com/jordyf15/thullo-api/user_boards/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func TestUserUsecase(t *testing.T) {
	suite.Run(t, new(userUsecaseSuite))
}

type userUsecaseSuite struct {
	suite.Suite

	usecase        user.Usecase
	userRepo       *ur.Repository
	tokenRepo      *tr.Repository
	oauthRepo      *or.Repository
	userBoardsRepo *ubr.Repository
	storage        *sr.Storage
}

func bcryptHash(str string) string {
	hashedStr, _ := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(hashedStr)
}

var (
	userID = primitive.NewObjectID()

	hashedPassword = bcryptHash("Password123!")
	user1          = &models.User{
		ID:                userID,
		Email:             "jojo@gmail.com",
		EncryptedPassword: hashedPassword,
		Username:          "jojo",
		Name:              "joseph joestar",
		Bio:               "a joestar",
		Images: []*models.Image{
			{URL: "image1", Width: 100},
			{URL: "image2", Width: 400},
		},
	}
)

func (s *userUsecaseSuite) SetupTest() {
	s.tokenRepo = new(tr.Repository)
	s.userRepo = new(ur.Repository)
	s.oauthRepo = new(or.Repository)
	s.userBoardsRepo = new(ubr.Repository)
	s.storage = new(sr.Storage)

	fieldExists := func(key, value string) bool {
		if key == "email" && value == "registered@gmail.com" {
			return true
		}
		if key == "username" && value == "alreadyexist" {
			return true
		}

		return false
	}

	s.userRepo.On("FieldExists", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(fieldExists, nil)
	s.userRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)
	s.userRepo.On("GetByEmail", mock.AnythingOfType("string")).Return(user1, nil)
	s.storage.On("UploadFile", mock.AnythingOfType("chan<- error"), mock.AnythingOfType("*sync.WaitGroup"), mock.AnythingOfType("*models.Image"), mock.AnythingOfType("*os.File"), mock.AnythingOfType("map[string]string")).Run(func(args mock.Arguments) {
		arg1 := args[0].(chan<- error)
		arg1 <- nil
		arg2 := args[1].(*sync.WaitGroup)
		arg2.Done()
	})
	s.storage.On("AssignImageURLToUser", mock.AnythingOfType("*models.User")).Return(nil)
	s.tokenRepo.On("Create", mock.AnythingOfType("*models.TokenSet")).Return(nil)
	s.userBoardsRepo.On("Create", mock.AnythingOfType("primitive.ObjectID")).Return(nil)

	s.usecase = usecase.NewUserUsecase(s.userRepo, s.tokenRepo, s.oauthRepo, s.userBoardsRepo, s.storage)
}

func (s *userUsecaseSuite) TestCreateInvalidFields() {
	user := &models.User{
		Email:    "",
		Name:     "",
		Username: "",
		Password: "",
	}

	result, err := s.usecase.Create(user, nil)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)

	expectedErrors := &custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrEmailAddressInvalid, custom_errors.ErrUsernameInvalid, custom_errors.ErrUsernameTooShort, custom_errors.ErrNameTooShort, custom_errors.ErrPasswordTooShort}}
	assert.Equal(s.T(), expectedErrors.Error(), err.Error())
}

func (s *userUsecaseSuite) TestCreateFieldUsernameAndEmailAlreadyExists() {
	user := &models.User{
		Email:    "registered@gmail.com",
		Name:     "joseph joestar",
		Username: "alreadyexist",
		Password: "Password123!",
	}

	result, err := s.usecase.Create(user, nil)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)

	expectedErrors := &custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrUsernameAlreadyExists, custom_errors.ErrEmailAddressAlreadyRegistered}}
	assert.Equal(s.T(), expectedErrors.Error(), err.Error())
}

func (s *userUsecaseSuite) TestLoginWrongPassword() {
	loginResponse, err := s.usecase.Login("jojo@gmail.com", "wrongPassword")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrCurrentPasswordWrong, err)
	assert.Nil(s.T(), loginResponse)
}

func (s *userUsecaseSuite) TestLoginSuccessful() {
	loginResponse, err := s.usecase.Login("jojo@gmail.com", "Password123!")

	assert.NoError(s.T(), err)

	data, isExist := loginResponse["data"].(*models.User)
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), userID, data.ID)
	assert.Equal(s.T(), "jojo@gmail.com", data.Email)
	assert.Equal(s.T(), "jojo", data.Username)
	assert.Equal(s.T(), "joseph joestar", data.Name)
	assert.Equal(s.T(), "a joestar", data.Bio)
	assert.Len(s.T(), data.Images, 2)
	assert.Equal(s.T(), "image1", data.Images[0].URL)
	assert.Equal(s.T(), uint(100), data.Images[0].Width)
	assert.Equal(s.T(), "image2", data.Images[1].URL)
	assert.Equal(s.T(), uint(400), data.Images[1].Width)
}
