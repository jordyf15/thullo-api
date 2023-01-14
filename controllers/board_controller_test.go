package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	usecaseMock.On("AddMember", mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("primitive.ObjectID")).Return(nil)
	usecaseMock.On("UpdateMemberRole", mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("string")).Return(nil)
	usecaseMock.On("DeleteMember", mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("primitive.ObjectID")).Return(nil)

	s.controller = controllers.NewBoardController(usecaseMock)
	s.response = httptest.NewRecorder()
	s.context, s.router = gin.CreateTestContext(s.response)

	s.router.POST("/boards", func(c *gin.Context) {
		c.Set("current_user_id", primitive.NewObjectID())
		c.Next()
	}, s.controller.Create)
	s.router.POST("/boards/:board_id/members", func(c *gin.Context) {
		c.Set("current_user_id", primitive.NewObjectID())
		c.Next()
	}, s.controller.AddMember)
	s.router.PATCH("/boards/:board_id/members/:member_id", func(c *gin.Context) {
		c.Set("current_user_id", primitive.NewObjectID())
		c.Next()
	}, s.controller.UpdateMemberRole)
	s.router.DELETE("/boards/:board_id/members/:member_id", func(c *gin.Context) {
		c.Set("current_user_id", primitive.NewObjectID())
		c.Next()
	}, s.controller.DeleteMember)
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

func (s *boardControllerSuite) TestAddMember() {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	memberID, _ := writer.CreateFormField("member_id")
	memberID.Write([]byte(primitive.NewObjectID().Hex()))
	writer.Close()

	s.context.Request, _ = http.NewRequest("POST", fmt.Sprintf("/boards/%s/members", primitive.NewObjectID().Hex()), buf)
	s.context.Request.Header.Set("Content-Type", writer.FormDataContentType())
	s.router.ServeHTTP(s.response, s.context.Request)

	assert.Equal(s.T(), http.StatusNoContent, s.response.Code)
}
func (s *boardControllerSuite) TestUpdateMemberRole() {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	role, _ := writer.CreateFormField("role")
	role.Write([]byte("admin"))
	writer.Close()

	s.context.Request, _ = http.NewRequest("PATCH", fmt.Sprintf("/boards/%s/members/%s", primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex()), buf)
	s.context.Request.Header.Set("Content-Type", writer.FormDataContentType())
	s.router.ServeHTTP(s.response, s.context.Request)

	assert.Equal(s.T(), http.StatusNoContent, s.response.Code)
}

func (s *boardControllerSuite) TestDeleteMember() {
	s.context.Request, _ = http.NewRequest("DELETE", fmt.Sprintf("/boards/%s/members/%s", primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex()), nil)

	s.router.ServeHTTP(s.response, s.context.Request)
}
