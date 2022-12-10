package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/jordyf15/thullo-api/models"
	"github.com/jordyf15/thullo-api/token"
)

type TokenController interface {
	RefreshAccessToken(c *gin.Context)
	DeleteRefreshToken(c *gin.Context)
}

type tokenController struct {
	usecase token.Usecase
}

func NewTokenController(usecase token.Usecase) TokenController {
	return &tokenController{usecase: usecase}
}

func (controller *tokenController) RefreshAccessToken(c *gin.Context) {
	refreshTokenStr := c.PostForm("refresh_token")
	refreshToken, err := parseRefreshToken(refreshTokenStr)

	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	newAccessToken, err := controller.usecase.Refresh(refreshToken)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"refresh_token": refreshToken.ToJWTString(),
		"access_token":  newAccessToken.ToJWTString(),
		"expires_at":    newAccessToken.ExpiresAt,
	})
}

func (controller *tokenController) DeleteRefreshToken(c *gin.Context) {
	refreshTokenStr := c.PostForm("refresh_token")
	refreshToken, err := parseRefreshToken(refreshTokenStr)

	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	err = controller.usecase.DeleteRefreshToken(refreshToken)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func parseRefreshToken(tokenStr string) (*models.RefreshToken, error) {
	refreshToken := &models.RefreshToken{}
	token, err := jwt.ParseWithClaims(tokenStr, refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_PASSWORD")), nil
	})

	if err != nil {
		return nil, custom_errors.ErrMalformedRefreshToken
	}

	if !token.Valid {
		return nil, custom_errors.ErrInvalidRefreshToken
	}

	return refreshToken, nil
}
