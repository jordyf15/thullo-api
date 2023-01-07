package controllers_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/thullo-api/board/mocks"
	"github.com/jordyf15/thullo-api/controllers"
	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBoardController(t *testing.T) {
	suite.Run(t, new(boardControllerSuite))
}

type boardControllerSuite struct {
	suite.Suite

	router     *gin.Engine
	controller controllers.BoardController
	response   *httptest.ResponseRecorder
	context    *gin.Context
}

func (s *boardControllerSuite) SetupTest() {
	usecaseMock := new(mocks.Usecase)

	usecaseMock.On("Create", mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(nil)

	s.controller = controllers.NewBoardController(usecaseMock)
	s.response = httptest.NewRecorder()
	s.context, s.router = gin.CreateTestContext(s.response)

	s.router.POST("/boards", func(c *gin.Context) {
		c.Set("current_user_id", primitive.NewObjectID())
		c.Next()
	}, s.controller.Create)
}

func (s *boardControllerSuite) TestCreateEmptyCover() {
	var receivedResponse map[string]interface{}
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	title, _ := writer.CreateFormField("title")
	title.Write([]byte("board 1"))
	visibility, _ := writer.CreateFormField("visibility")
	visibility.Write([]byte("public"))
	writer.Close()

	s.context.Request, _ = http.NewRequest("POST", "/boards", buf)
	s.context.Request.Header.Set("Content-Type", writer.FormDataContentType())
	s.router.ServeHTTP(s.response, s.context.Request)
	json.NewDecoder(s.response.Body).Decode(&receivedResponse)

	assert.Equal(s.T(), http.StatusBadRequest, s.response.Code)

	errors, isExist := receivedResponse["errors"].([]interface{})
	assert.True(s.T(), isExist)

	error1 := errors[0].(map[string]interface{})

	msg, isExist := error1["message"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), custom_errors.ErrBoardCoverEmpty.Message, msg)

	code, isExist := error1["code"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), float64(custom_errors.ErrBoardCoverEmpty.Code), code)
}

func (s *boardControllerSuite) TestCreateMalformedCover() {
	var receivedResponse map[string]interface{}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	title, _ := writer.CreateFormField("title")
	title.Write([]byte("board 1"))
	visibility, _ := writer.CreateFormField("visibility")
	visibility.Write([]byte("public"))
	cover, _ := writer.CreateFormField("cover")
	cover.Write([]byte("malformedCover"))
	writer.Close()

	s.context.Request, _ = http.NewRequest("POST", "/boards", buf)
	s.context.Request.Header.Set("Content-Type", writer.FormDataContentType())
	s.router.ServeHTTP(s.response, s.context.Request)
	json.NewDecoder(s.response.Body).Decode(&receivedResponse)

	assert.Equal(s.T(), http.StatusBadRequest, s.response.Code)

	errors, isExist := receivedResponse["errors"].([]interface{})
	assert.True(s.T(), isExist)

	error1 := errors[0].(map[string]interface{})

	msg, isExist := error1["message"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), custom_errors.ErrMalformedCover.Message, msg)

	code, isExist := error1["code"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), float64(custom_errors.ErrMalformedCover.Code), code)
}

func (s *boardControllerSuite) TestCreate() {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	title, _ := writer.CreateFormField("title")
	title.Write([]byte("board 1"))
	visibility, _ := writer.CreateFormField("visibility")
	visibility.Write([]byte("public"))
	cover, _ := writer.CreateFormField("cover")
	cover.Write([]byte("unsplash:unsplashid1:0.5"))
	writer.Close()

	s.context.Request, _ = http.NewRequest("POST", "/boards", buf)
	s.context.Request.Header.Set("Content-Type", writer.FormDataContentType())
	s.router.ServeHTTP(s.response, s.context.Request)

	assert.Equal(s.T(), http.StatusNoContent, s.response.Code)
}
