package usecase

import (
	"os"
	"time"

	"github.com/jordyf15/thullo-api/models"
	"github.com/jordyf15/thullo-api/token"
	"github.com/jordyf15/thullo-api/utils"
)

type tokenUsecase struct {
	repo token.Repository
}

func NewTokenUsecase(repo token.Repository) token.Usecase {
	return &tokenUsecase{repo: repo}
}

func (usecase *tokenUsecase) Refresh(refreshToken *models.RefreshToken) (*models.AccessToken, error) {
	hashedRefreshTokenID := utils.ToSHA256(refreshToken.Id)
	tokenSet, err := usecase.repo.GetTokenSet(refreshToken.UserID, hashedRefreshTokenID, true)
	if err != nil {
		return nil, err
	}

	refreshToken.Id = utils.RandString(8)
	tokenSet.PrevRefreshTokenID = &hashedRefreshTokenID
	tokenSet.RefreshTokenID = utils.ToSHA256(refreshToken.Id)
	err = usecase.repo.Update(tokenSet)

	if err != nil {
		return nil, err
	}

	var expiration time.Time
	if os.Getenv("ENV") == "development" {
		expiration = time.Now().Add(time.Hour * 24)
	} else {
		expiration = time.Now().Add(time.Hour * 1)
	}
	accessToken := (&models.AccessToken{UserID: tokenSet.UserID, RefreshTokenID: tokenSet.RefreshTokenID}).SetExpiration(expiration)
	accessToken.Id = utils.RandString(8)
	usecase.repo.Save(accessToken)

	return accessToken, nil
}

func (usecase *tokenUsecase) Use(token *models.AccessToken) error {
	if usecase.repo.Exists(token) {
		tokenSet, err := usecase.repo.GetTokenSet(token.UserID, token.RefreshTokenID, false)
		if err != nil {
			return err
		}

		usecase.repo.Updates(tokenSet, map[string]interface{}{"prt_id": nil})
		usecase.repo.Remove(token)
	}
	return nil
}

func (usecase *tokenUsecase) DeleteRefreshToken(token *models.RefreshToken) error {
	return usecase.repo.Delete(&models.TokenSet{UserID: token.UserID, RefreshTokenID: utils.ToSHA256(token.Id)})
}
