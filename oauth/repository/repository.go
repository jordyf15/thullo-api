package repository

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jordyf15/thullo-api/models"
	"github.com/jordyf15/thullo-api/oauth"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type oauthRepository struct {
	httpclient *http.Client
}

func NewOauthRepository(httpClient *http.Client) oauth.Repository {
	return &oauthRepository{httpclient: httpClient}
}

func (repo *oauthRepository) GetGoogleTokenInfo(token string) (*models.GoogleTokenInfo, error) {
	ctx := context.Background()

	oauth2Service, err := oauth2.NewService(ctx, option.WithHTTPClient(repo.httpclient))
	if err != nil {
		return nil, err
	}

	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(token)
	_, err = tokenInfoCall.Do()
	if err != nil {
		return nil, err
	}

	tokenInfo := &models.GoogleTokenInfo{}
	_, _ = jwt.ParseWithClaims(token, tokenInfo, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})

	return tokenInfo, nil
}
