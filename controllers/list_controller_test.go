package controllers_test

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/thullo-api/controllers"
	"github.com/jordyf15/thullo-api/list/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestListController(t *testing.T) {
	suite.Run(t, new(listControllerSuite))
}

type listControllerSuite struct {
	suite.Suite

	router     *gin.Engine
	controller controllers.ListController
	response   *httptest.ResponseRecorder
	context    *gin.Context
	usecase    *mocks.Usecase
}

func (s *listControllerSuite) SetupTest() {
	s.usecase = new(mocks.Usecase)

	s.usecase.On("Create", mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("string")).Return(nil)
	s.usecase.On("UpdateTitle", mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("string")).Return(nil)
	s.usecase.On("UpdatePosition", mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("int")).Return(nil)

	s.controller = controllers.NewListController(s.usecase)
	s.response = httptest.NewRecorder()
	s.context, s.router = gin.CreateTestContext(s.response)

	s.router.POST("/boards/:board_id/lists", func(c *gin.Context) {
		c.Set("current_user_id", primitive.NewObjectID())
		c.Next()
	}, s.controller.Create)
	s.router.PATCH("/boards/:board_id/lists/:list_id", func(c *gin.Context) {
		c.Set("current_user_id", primitive.NewObjectID())
		c.Next()
	}, s.controller.Update)
}

func (s *listControllerSuite) TestCreate() {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	title, _ := writer.CreateFormField("title")
	title.Write([]byte("list 1"))
	writer.Close()

	s.context.Request, _ = http.NewRequest("POST", fmt.Sprintf("/boards/%s/lists", primitive.NewObjectID().Hex()), buf)
	s.context.Request.Header.Set("Content-Type", writer.FormDataContentType())
	s.router.ServeHTTP(s.response, s.context.Request)

	assert.Equal(s.T(), http.StatusNoContent, s.response.Code)
}

func (s *listControllerSuite) TestUpdateTitle() {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	title, _ := writer.CreateFormField("title")
	title.Write([]byte("list 1"))
	writer.Close()

	s.context.Request, _ = http.NewRequest("PATCH", fmt.Sprintf("/boards/%s/lists/%s", primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex()), buf)
	s.context.Request.Header.Set("Content-Type", writer.FormDataContentType())
	s.router.ServeHTTP(s.response, s.context.Request)

	assert.Equal(s.T(), http.StatusNoContent, s.response.Code)
	s.usecase.AssertNumberOfCalls(s.T(), "UpdateTitle", 1)
	s.usecase.AssertNumberOfCalls(s.T(), "UpdatePosition", 0)
}

func (s *listControllerSuite) TestUpdatePosition() {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	position, _ := writer.CreateFormField("position")
	position.Write([]byte("10"))
	writer.Close()

	s.context.Request, _ = http.NewRequest("PATCH", fmt.Sprintf("/boards/%s/lists/%s", primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex()), buf)
	s.context.Request.Header.Set("Content-Type", writer.FormDataContentType())
	s.router.ServeHTTP(s.response, s.context.Request)

	assert.Equal(s.T(), http.StatusNoContent, s.response.Code)
	s.usecase.AssertNumberOfCalls(s.T(), "UpdateTitle", 0)
	s.usecase.AssertNumberOfCalls(s.T(), "UpdatePosition", 1)
}
