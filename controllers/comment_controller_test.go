package controllers_test

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/thullo-api/comment/mocks"
	"github.com/jordyf15/thullo-api/controllers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCommentController(t *testing.T) {
	suite.Run(t, new(commentControllerSuite))
}

type commentControllerSuite struct {
	suite.Suite

	router     *gin.Engine
	controller controllers.CommentController
	response   *httptest.ResponseRecorder
	context    *gin.Context
	usecase    *mocks.Usecase
}

func (s *commentControllerSuite) SetupTest() {
	s.usecase = new(mocks.Usecase)

	s.usecase.On("Create", mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("string")).Return(nil)
	s.usecase.On("Update", mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("string")).Return(nil)

	s.controller = controllers.NewCommentController(s.usecase)
	s.response = httptest.NewRecorder()
	s.context, s.router = gin.CreateTestContext(s.response)

	s.router.POST("/boards/:board_id/lists/:list_id/cards/:card_id/comments", func(c *gin.Context) {
		c.Set("current_user_id", primitive.NewObjectID())
		c.Next()
	}, s.controller.Create)
	s.router.PATCH("/boards/:board_id/lists/:list_id/cards/:card_id/comments/:comment_id", func(c *gin.Context) {
		c.Set("current_user_id", primitive.NewObjectID())
		c.Next()
	}, s.controller.Update)
}

func (s *commentControllerSuite) TestCreate() {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	comment, _ := writer.CreateFormField("comment")
	comment.Write([]byte("comment 1"))
	writer.Close()

	s.context.Request, _ = http.NewRequest("POST", fmt.Sprintf("/boards/%s/lists/%s/cards/%s/comments", primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex()), buf)
	s.context.Request.Header.Set("Content-Type", writer.FormDataContentType())
	s.router.ServeHTTP(s.response, s.context.Request)

	assert.Equal(s.T(), http.StatusNoContent, s.response.Code)
}

func (s *commentControllerSuite) TestUpdate() {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	comment, _ := writer.CreateFormField("comment")
	comment.Write([]byte("updated comment"))
	writer.Close()

	s.context.Request, _ = http.NewRequest("PATCH", fmt.Sprintf("/boards/%s/lists/%s/cards/%s/comments/%s", primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex()), buf)
	s.context.Request.Header.Set("Content-Type", writer.FormDataContentType())
	s.router.ServeHTTP(s.response, s.context.Request)

	assert.Equal(s.T(), http.StatusNoContent, s.response.Code)
}
