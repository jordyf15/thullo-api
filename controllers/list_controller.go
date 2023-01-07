package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/thullo-api/list"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListController interface {
	Create(c *gin.Context)
}

type listController struct {
	usecase list.Usecase
}

func NewListController(usecase list.Usecase) ListController {
	return &listController{usecase: usecase}
}

func (controller *listController) Create(c *gin.Context) {
	boardIDStr := c.Param("board_id")
	title := c.PostForm("title")

	boardID, err := primitive.ObjectIDFromHex(boardIDStr)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	err = controller.usecase.Create(boardID, title)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
