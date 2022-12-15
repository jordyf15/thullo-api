package user

import (
	"github.com/jordyf15/thullo-api/models"
)

var (
	DisplayPictureSizes = []uint{100, 400}
)

type Repository interface {
	Create(user *models.User) error
	FieldExists(key string, value string) (bool, error)
}

type Usecase interface {
	Create(*models.User) (map[string]interface{}, error)
	For(user *models.User) InstanceUsecase
}

type InstanceUsecase interface {
	GenerateTokens() (*models.AccessToken, *models.RefreshToken, error)
}
