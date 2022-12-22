package controllers_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/thullo-api/controllers"
	"github.com/jordyf15/thullo-api/models"
	"github.com/jordyf15/thullo-api/user/mocks"
	"github.com/jordyf15/thullo-api/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserController(t *testing.T) {
	suite.Run(t, new(userControllerSuite))
}

type userControllerSuite struct {
	suite.Suite

	router     *gin.Engine
	controller controllers.UserController
	response   *httptest.ResponseRecorder
	context    *gin.Context
}

var (
	uscUserID = primitive.NewObjectID()

	uscImage1 = &models.Image{
		Width: 100,
		URL:   "http://storage/image1.jpg",
	}
	uscImage2 = &models.Image{
		Width: 400,
		URL:   "http://storage/image2.jpg",
	}
	uscImages = models.Images{
		uscImage1,
		uscImage2,
	}

	uscUser = &models.User{
		ID:        uscUserID,
		Email:     "jojo@gmail.com",
		Username:  "jojo",
		Name:      "joseph joestar",
		Bio:       "a joestar",
		Images:    uscImages,
		CreatedAt: time.Now(),
	}

	uscLoginResponse = utils.DataResponse(uscUser, map[string]interface{}{
		"access_token":  "accessToken",
		"refresh_token": "refreshToken",
		"expires_at":    1,
	})
)

func (s *userControllerSuite) SetupTest() {
	usecaseMock := new(mocks.Usecase)

	usecaseMock.On("LoginWithGoogle", mock.AnythingOfType("string")).Return(uscLoginResponse, nil)
	usecaseMock.On("Login", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(uscLoginResponse, nil)

	s.controller = controllers.NewUserController(usecaseMock)
	s.response = httptest.NewRecorder()
	s.context, s.router = gin.CreateTestContext(s.response)

	s.router.POST("/login/google", s.controller.LoginWithGoogle)
	s.router.POST("/login", s.controller.Login)
}

func (s *userControllerSuite) TestLogin() {
	var receivedResponse map[string]interface{}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	emailWr, _ := writer.CreateFormField("email")
	emailWr.Write([]byte("jojo@gmail.com"))
	password, _ := writer.CreateFormField("password")
	password.Write([]byte("Password123!"))
	writer.Close()

	s.context.Request, _ = http.NewRequest("POST", "/login", buf)
	s.context.Request.Header.Set("Content-Type", writer.FormDataContentType())
	s.router.ServeHTTP(s.response, s.context.Request)
	json.NewDecoder(s.response.Body).Decode(&receivedResponse)

	assert.Equal(s.T(), http.StatusOK, s.response.Code)

	data, isExist := receivedResponse["data"].(map[string]interface{})
	assert.True(s.T(), isExist)

	id, isExist := data["id"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), uscUserID.Hex(), id)

	email, isExist := data["email"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), "jojo@gmail.com", email)

	username, isExist := data["username"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), "jojo", username)

	name, isExist := data["name"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), "joseph joestar", name)

	bio, isExist := data["bio"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), "a joestar", bio)

	images, isExist := data["images"].([]interface{})
	assert.True(s.T(), isExist)

	image1 := images[0].(map[string]interface{})
	width, isExist := image1["width"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), float64(100), width)
	url, isExist := image1["url"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), "http://storage/image1.jpg", url)

	image2 := images[1].(map[string]interface{})
	width, isExist = image2["width"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), float64(400), width)
	url, isExist = image2["url"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), "http://storage/image2.jpg", url)

	_, isExist = data["created_at"]
	assert.True(s.T(), isExist)

	meta, isExist := receivedResponse["meta"].(map[string]interface{})
	assert.True(s.T(), isExist)

	accessToken, isExist := meta["access_token"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), "accessToken", accessToken)

	refreshToken, isExist := meta["refresh_token"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), "refreshToken", refreshToken)

	expiresAt, isExist := meta["expires_at"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), float64(1), expiresAt)
}

func (s *userControllerSuite) TestLoginWithGoogle() {
	var receivedResponse map[string]interface{}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	token, _ := writer.CreateFormField("token")
	token.Write([]byte("token"))
	writer.Close()

	s.context.Request, _ = http.NewRequest("POST", "/login/google", buf)
	s.context.Request.Header.Set("Content-Type", writer.FormDataContentType())
	s.router.ServeHTTP(s.response, s.context.Request)
	json.NewDecoder(s.response.Body).Decode(&receivedResponse)

	assert.Equal(s.T(), http.StatusOK, s.response.Code)

	data, isExist := receivedResponse["data"].(map[string]interface{})
	assert.True(s.T(), isExist)

	id, isExist := data["id"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), uscUserID.Hex(), id)

	email, isExist := data["email"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), "jojo@gmail.com", email)

	username, isExist := data["username"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), "jojo", username)

	name, isExist := data["name"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), "joseph joestar", name)

	bio, isExist := data["bio"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), "a joestar", bio)

	images, isExist := data["images"].([]interface{})
	assert.True(s.T(), isExist)

	image1 := images[0].(map[string]interface{})
	width, isExist := image1["width"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), float64(100), width)
	url, isExist := image1["url"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), "http://storage/image1.jpg", url)

	image2 := images[1].(map[string]interface{})
	width, isExist = image2["width"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), float64(400), width)
	url, isExist = image2["url"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), "http://storage/image2.jpg", url)

	_, isExist = data["created_at"]
	assert.True(s.T(), isExist)

	meta, isExist := receivedResponse["meta"].(map[string]interface{})
	assert.True(s.T(), isExist)

	accessToken, isExist := meta["access_token"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), "accessToken", accessToken)

	refreshToken, isExist := meta["refresh_token"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), "refreshToken", refreshToken)

	expiresAt, isExist := meta["expires_at"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), float64(1), expiresAt)
}
