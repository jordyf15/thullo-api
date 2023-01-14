package controllers_test

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/thullo-api/card/mocks"
	"github.com/jordyf15/thullo-api/controllers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCardController(t *testing.T) {
	suite.Run(t, new(cardControllerSuite))
}

type cardControllerSuite struct {
	suite.Suite

	router     *gin.Engine
	controller controllers.CardController
	response   *httptest.ResponseRecorder
	context    *gin.Context
	usecase    *mocks.Usecase
}

func (s *cardControllerSuite) SetupTest() {
	s.usecase = new(mocks.Usecase)

	s.usecase.On("Create", mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("string")).Return(nil)

	s.controller = controllers.NewCardController(s.usecase)
	s.response = httptest.NewRecorder()
	s.context, s.router = gin.CreateTestContext(s.response)

	s.router.POST("/boards/:board_id/lists/:list_id/cards", func(c *gin.Context) {
		c.Set("current_user_id", primitive.NewObjectID())
		c.Next()
	}, s.controller.Create)
}

func (s *cardControllerSuite) TestCreate() {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	title, _ := writer.CreateFormField("title")
	title.Write([]byte("card 1"))
	writer.Close()

	s.context.Request, _ = http.NewRequest("POST", fmt.Sprintf("/boards/%s/lists/%s/cards", primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex()), buf)
	s.context.Request.Header.Set("Content-Type", writer.FormDataContentType())
	s.router.ServeHTTP(s.response, s.context.Request)

	assert.Equal(s.T(), http.StatusNoContent, s.response.Code)
}
