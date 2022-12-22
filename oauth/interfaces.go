package oauth

import "github.com/jordyf15/thullo-api/models"

type Repository interface {
	GetGoogleTokenInfo(token string) (*models.GoogleTokenInfo, error)
}
