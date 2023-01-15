package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/thullo-api/comment"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentController interface {
	Create(c *gin.Context)
}

type commentController struct {
	usecase comment.Usecase
}

func NewCommentController(usecase comment.Usecase) CommentController {
	return &commentController{usecase: usecase}
}

func (controller *commentController) Create(c *gin.Context) {
	requesterID := c.MustGet("current_user_id").(primitive.ObjectID)
	boardIDStr := c.Param("board_id")
	cardIDStr := c.Param("card_id")
	comment := strings.TrimSpace(c.PostForm("comment"))

	boardID, err := primitive.ObjectIDFromHex(boardIDStr)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	cardID, err := primitive.ObjectIDFromHex(cardIDStr)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	err = controller.usecase.Create(requesterID, boardID, cardID, comment)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
