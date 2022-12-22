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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestUserUsecase(t *testing.T) {
	suite.Run(t, new(userUsecaseSuite))
}

type userUsecaseSuite struct {
	suite.Suite

	usecase   user.Usecase
	userRepo  *ur.Repository
	tokenRepo *tr.Repository
	oauthRepo *or.Repository
	storage   *sr.Storage
}

func (s *userUsecaseSuite) SetupTest() {
	s.tokenRepo = new(tr.Repository)
	s.userRepo = new(ur.Repository)
	s.oauthRepo = new(or.Repository)
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
	s.storage.On("UploadFile", mock.AnythingOfType("chan<- error"), mock.AnythingOfType("*sync.WaitGroup"), mock.AnythingOfType("*models.Image"), mock.AnythingOfType("*os.File"), mock.AnythingOfType("map[string]string")).Run(func(args mock.Arguments) {
		arg1 := args[0].(chan<- error)
		arg1 <- nil
		arg2 := args[1].(*sync.WaitGroup)
		arg2.Done()
	})
	s.userRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	s.usecase = usecase.NewUserUsecase(s.userRepo, s.tokenRepo, s.oauthRepo, s.storage)
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
