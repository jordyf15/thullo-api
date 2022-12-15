package usecase

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/jordyf15/thullo-api/models"
	"github.com/jordyf15/thullo-api/storage"
	"github.com/jordyf15/thullo-api/token"
	"github.com/jordyf15/thullo-api/user"
	"github.com/jordyf15/thullo-api/utils"
)

type userUsecase struct {
	userRepo  user.Repository
	tokenRepo token.Repository
	storage   storage.Storage
}

type userInstanceUsecase struct {
	user *models.User
	userUsecase
}

func NewUserUsecase(userRepo user.Repository, tokenRepo token.Repository, storage storage.Storage) user.Usecase {
	return &userUsecase{userRepo: userRepo, tokenRepo: tokenRepo, storage: storage}
}

func (usecase *userUsecase) Create(_user *models.User) (map[string]interface{}, error) {
	var err error
	errors := make([]error, 0)

	validateFieldErrors := _user.VerifyFields()
	if len(validateFieldErrors) > 0 {
		errors = append(errors, validateFieldErrors...)
	}

	isUsernameExist, err := usecase.userRepo.FieldExists("username", _user.Username)
	if err != nil {
		return nil, err
	}
	if isUsernameExist {
		errors = append(errors, custom_errors.ErrUsernameAlreadyExists)
	}

	isEmailExist, err := usecase.userRepo.FieldExists("email", _user.Email)
	if err != nil {
		return nil, err
	}
	if isEmailExist {
		errors = append(errors, custom_errors.ErrEmailAddressAlreadyRegistered)
	}

	err = _user.SetPassword(_user.Password)
	if err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return nil, &custom_errors.MultipleErrors{Errors: errors}
	}

	avatar, err := utils.GenerateAvatar(_user.Initials(), 800, 400)
	if err != nil {
		return nil, err
	}

	defer os.Remove(avatar.Name())

	userAvatar := utils.NewNamedFileReader(avatar, avatar.Name())

	_user.Images = make([]*models.Image, len(user.DisplayPictureSizes))
	for i, width := range user.DisplayPictureSizes {
		image := &models.Image{}
		image.Width = width

		_user.Images[i] = image
	}

	uploadChannels := make(chan error, len(_user.Images))
	var wg sync.WaitGroup

	wg.Add(len(_user.Images))

	for _, img := range _user.Images {
		var resizedImageFile *os.File
		resizedImageFile, err = utils.ResizeImage(userAvatar, int(img.Width))
		if err != nil {
			break
		}

		defer os.Remove(resizedImageFile.Name())

		name := utils.RandString(8)
		fileName := fmt.Sprintf("%s.%s", name, utils.GetFileExtension(resizedImageFile.Name()))

		metaData := map[string]string{
			"name":        fileName,
			"title":       name,
			"description": fmt.Sprintf("profile picture of %s with width and height of %v", _user.Username, img.Width),
		}

		go usecase.storage.UploadFile(uploadChannels, &wg, img, resizedImageFile, metaData)
	}

	wg.Wait()
	close(uploadChannels)

	for err = range uploadChannels {
		if err != nil {
			fmt.Println(err)
		}
	}

	<-uploadChannels

	if err := usecase.userRepo.Create(_user); err != nil {
		return nil, err
	}

	accessToken, refreshToken, _ := usecase.For(_user).GenerateTokens()

	response := utils.DataResponse(_user, map[string]interface{}{
		"access_token":  accessToken.ToJWTString(),
		"refresh_token": refreshToken.ToJWTString(),
		"expires_at":    accessToken.ExpiresAt,
	})

	return response, nil
}

func (usecase *userUsecase) For(user *models.User) user.InstanceUsecase {
	instanceUsecase := &userInstanceUsecase{user: user, userUsecase: *usecase}
	return instanceUsecase
}

// userInstanceUsecase
func (usecase *userInstanceUsecase) GenerateTokens() (*models.AccessToken, *models.RefreshToken, error) {
	refreshToken := (&models.RefreshToken{UserID: usecase.user.ID})
	refreshToken.Id = utils.RandString(8)

	accessToken := (&models.AccessToken{UserID: usecase.user.ID}).SetExpiration(time.Now().Add(time.Hour * 1))
	accessToken.RefreshTokenID = utils.ToSHA256(refreshToken.Id)

	tokenSet := &models.TokenSet{UserID: usecase.user.ID, RefreshTokenID: accessToken.RefreshTokenID}
	err := usecase.tokenRepo.Create(tokenSet)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}
