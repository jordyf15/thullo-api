package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jordyf15/thullo-api/controllers"
	"github.com/jordyf15/thullo-api/models"
	"github.com/jordyf15/thullo-api/token/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTokenController(t *testing.T) {
	suite.Run(t, new(tokenControllerSuite))
}

type tokenControllerSuite struct {
	suite.Suite
	router     *gin.Engine
	controller controllers.TokenController
	response   *httptest.ResponseRecorder
	context    *gin.Context
}

func (s *tokenControllerSuite) SetupTest() {
	usecaseMock := new(mocks.Usecase)
	accessToken := &models.AccessToken{}
	usecaseMock.On("Refresh", mock.AnythingOfType("*models.RefreshToken")).Return(accessToken, nil)
	usecaseMock.On("DeleteRefreshToken", mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	s.controller = controllers.NewTokenController(usecaseMock)
	s.response = httptest.NewRecorder()
	s.context, s.router = gin.CreateTestContext(s.response)
	s.router.POST("/tokens/refresh", s.controller.RefreshAccessToken)
	s.router.DELETE("/tokens/remove", s.controller.DeleteRefreshToken)
}

func (s *tokenControllerSuite) TestRefreshAccessToken() {
	var receivedResponse map[string]interface{}
	refreshToken := models.RefreshToken{
		UserID: primitive.NewObjectID(),
		Token: models.Token{StandardClaims: jwt.StandardClaims{
			Id:        "123456789",
			ExpiresAt: time.Now().Add(time.Hour * 5).Unix(),
		}},
	}
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	rt, _ := writer.CreateFormField("refresh_token")
	rt.Write([]byte(refreshToken.ToJWTString()))
	writer.Close()

	s.context.Request, _ = http.NewRequest("POST", "/tokens/refresh", buf)
	s.context.Request.Header.Set("Content-Type", writer.FormDataContentType())
	s.router.ServeHTTP(s.response, s.context.Request)
	json.NewDecoder(s.response.Body).Decode(&receivedResponse)

	assert.Equal(s.T(), http.StatusOK, s.response.Code)
	refreshTokenResp, isExist := receivedResponse["refresh_token"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), "string", fmt.Sprintf("%T", refreshTokenResp))
	accessToken, isExist := receivedResponse["access_token"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), "string", fmt.Sprintf("%T", accessToken))
	expiresAt, isExist := receivedResponse["expires_at"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), "float64", fmt.Sprintf("%T", expiresAt))
}

func (s *tokenControllerSuite) TestDeleteRefreshToken() {
	refreshToken := models.RefreshToken{
		UserID: primitive.NewObjectID(),
		Token: models.Token{StandardClaims: jwt.StandardClaims{
			Id:        "123456789",
			ExpiresAt: time.Now().Add(time.Hour * 5).Unix(),
		}},
	}
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	rt, _ := writer.CreateFormField("refresh_token")
	rt.Write([]byte(refreshToken.ToJWTString()))
	writer.Close()
	s.context.Request, _ = http.NewRequest("DELETE", "/tokens/remove", buf)
	s.context.Request.Header.Set("Content-Type", writer.FormDataContentType())
	s.router.ServeHTTP(s.response, s.context.Request)

	assert.Equal(s.T(), http.StatusNoContent, s.response.Code)
}
