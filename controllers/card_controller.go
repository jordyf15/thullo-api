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
	listIDStr := c.Param("list_id")

	listID, err := primitive.ObjectIDFromHex(listIDStr)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	title := strings.TrimSpace(c.PostForm("title"))

	err = controller.usecase.Create(listID, title)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
