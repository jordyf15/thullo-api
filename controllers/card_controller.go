package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/thullo-api/card"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CardController interface {
	Create(c *gin.Context)
}

type cardController struct {
	usecase card.Usecase
}

func NewCardController(usecase card.Usecase) CardController {
	return &cardController{usecase: usecase}
}

func (controller *cardController) Create(c *gin.Context) {
	userID := c.MustGet("current_user_id").(primitive.ObjectID)
	boardIDStr := c.Param("board_id")
	listIDStr := c.Param("list_id")

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

	title := strings.TrimSpace(c.PostForm("title"))

	err = controller.usecase.Create(userID, boardID, listID, title)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
