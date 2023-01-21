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
	Update(c *gin.Context)
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
	listIDStr := c.Param("list_id")
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

	listID, err := primitive.ObjectIDFromHex(listIDStr)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	err = controller.usecase.Create(requesterID, boardID, listID, cardID, comment)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (controller *commentController) Update(c *gin.Context) {
	requesterID := c.MustGet("current_user_id").(primitive.ObjectID)
	boardIDStr := c.Param("board_id")
	cardIDStr := c.Param("card_id")
	listIDStr := c.Param("list_id")
	commentIDStr := c.Param("comment_id")
	comment := strings.TrimSpace(c.PostForm("comment"))

	boardID, err := primitive.ObjectIDFromHex(boardIDStr)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	listID, err := primitive.ObjectIDFromHex(listIDStr)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	cardID, err := primitive.ObjectIDFromHex(cardIDStr)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	commentID, err := primitive.ObjectIDFromHex(commentIDStr)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	err = controller.usecase.Update(requesterID, boardID, listID, cardID, commentID, comment)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
