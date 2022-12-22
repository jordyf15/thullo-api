package user

import (
	"github.com/jordyf15/thullo-api/models"
	"github.com/jordyf15/thullo-api/utils"
)

var (
	DisplayPictureSizes = []uint{100, 400}
)

type Repository interface {
	Create(user *models.User) error
	FieldExists(key string, value string) (bool, error)
	GetByEmail(email string) (*models.User, error)
}

type Usecase interface {
	Create(user *models.User, imageFile utils.NamedFileReader) (map[string]interface{}, error)
	For(user *models.User) InstanceUsecase
	LoginWithGoogle(token string) (map[string]interface{}, error)
	Login(email, password string) (map[string]interface{}, error)
}

type InstanceUsecase interface {
	GenerateTokens() (*models.AccessToken, *models.RefreshToken, error)
}
