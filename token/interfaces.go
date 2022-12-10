package token

import (
	"github.com/jordyf15/thullo-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DefaultTokenLimitPerUser = 5
)

type Repository interface {
	GetTokenSet(userID primitive.ObjectID, hashedRefreshTokenID string, includeParent bool) (*models.TokenSet, error)
	Save(accessToken *models.AccessToken) error
	Exists(accessToken *models.AccessToken) bool
	Remove(accessToken *models.AccessToken) error
	Create(tokenSet *models.TokenSet) error
	Update(tokenSet *models.TokenSet) error
	Updates(tokenSet *models.TokenSet, changes map[string]interface{}) error
	Delete(tokenSet *models.TokenSet) error
}

type Usecase interface {
	Refresh(token *models.RefreshToken) (*models.AccessToken, error)
	Use(token *models.AccessToken) error
	DeleteRefreshToken(token *models.RefreshToken) error
}
