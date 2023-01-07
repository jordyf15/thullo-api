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
}

func (s *listControllerSuite) SetupTest() {
	usecaseMock := new(mocks.Usecase)

	usecaseMock.On("Create", mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("string")).Return(nil)

	s.controller = controllers.NewListController(usecaseMock)
	s.response = httptest.NewRecorder()
	s.context, s.router = gin.CreateTestContext(s.response)

	s.router.POST("/boards/:board_id/lists", s.controller.Create)
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
