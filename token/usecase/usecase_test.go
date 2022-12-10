package usecase_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jordyf15/thullo-api/models"
	"github.com/jordyf15/thullo-api/token"
	"github.com/jordyf15/thullo-api/token/mocks"
	"github.com/jordyf15/thullo-api/token/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTokenUsecase(t *testing.T) {
	suite.Run(t, new(tokenUsecaseSuite))
}

type tokenUsecaseSuite struct {
	suite.Suite
	usecase   token.Usecase
	tokenRepo *mocks.Repository
}

var (
	tokenID  = primitive.NewObjectID()
	userID   = primitive.NewObjectID()
	tokenSet = &models.TokenSet{
		ID:                 tokenID,
		UserID:             userID,
		RefreshTokenID:     "refreshTokenId",
		UpdatedAt:          time.Now(),
		PrevRefreshTokenID: nil,
	}
)

func (s *tokenUsecaseSuite) SetupTest() {
	s.tokenRepo = new(mocks.Repository)

	s.tokenRepo.On("GetTokenSet", mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("string"), mock.AnythingOfType("bool")).Return(tokenSet, nil)
	s.tokenRepo.On("Update", mock.AnythingOfType("*models.TokenSet")).Return(nil)
	s.tokenRepo.On("Save", mock.AnythingOfType("*models.AccessToken")).Return(nil)
	s.tokenRepo.On("Exists", mock.AnythingOfType("*models.AccessToken")).Return(true)
	s.tokenRepo.On("Updates", mock.AnythingOfType("*models.TokenSet"), mock.AnythingOfType("map[string]interface {}")).Return(nil)
	s.tokenRepo.On("Remove", mock.AnythingOfType("*models.AccessToken")).Return(nil)
	s.tokenRepo.On("Delete", mock.AnythingOfType("*models.TokenSet")).Return(nil)
	s.usecase = usecase.NewTokenUsecase(s.tokenRepo)
}

func (s *tokenUsecaseSuite) TestRefresh() {
	refreshToken := &models.RefreshToken{
		UserID: userID,
	}
	accessToken, err := s.usecase.Refresh(refreshToken)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), accessToken.RefreshTokenID)
	assert.Equal(s.T(), "string", fmt.Sprintf("%T", accessToken.RefreshTokenID))
	assert.Equal(s.T(), userID, accessToken.UserID)
	assert.NotEmpty(s.T(), accessToken.Id)
}

func (s *tokenUsecaseSuite) TestUse() {
	accessToken := &models.AccessToken{
		UserID:         userID,
		RefreshTokenID: "refreshTokenId",
	}

	err := s.usecase.Use(accessToken)

	assert.NoError(s.T(), err)
}

func (s *tokenUsecaseSuite) TestDeleteRefreshToken() {
	refreshToken := &models.RefreshToken{
		UserID: userID,
	}

	err := s.usecase.DeleteRefreshToken(refreshToken)

	assert.NoError(s.T(), err)
}
