package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/thullo-api/list"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListController interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
}

type listController struct {
	usecase list.Usecase
}

func NewListController(usecase list.Usecase) ListController {
	return &listController{usecase: usecase}
}

func (controller *listController) Create(c *gin.Context) {
	userID := c.MustGet("current_user_id").(primitive.ObjectID)
	boardIDStr := c.Param("board_id")
	title := strings.TrimSpace(c.PostForm("title"))

	boardID, err := primitive.ObjectIDFromHex(boardIDStr)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	err = controller.usecase.Create(userID, boardID, title)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (controller *listController) Update(c *gin.Context) {
	userID := c.MustGet("current_user_id").(primitive.ObjectID)
	listIDStr := c.Param("list_id")
	boardIDStr := c.Param("board_id")

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

	title, isExist := c.GetPostForm("title")
	if isExist {
		err := controller.usecase.UpdateTitle(userID, boardID, listID, strings.TrimSpace(title))
		if err != nil {
			respondBasedOnError(c, err)
			return
		}
	}

	positionStr, isExist := c.GetPostForm("position")
	if isExist {
		position, err := strconv.Atoi(positionStr)
		if err != nil {
			respondBasedOnError(c, err)
			return
		}

		err = controller.usecase.UpdatePosition(userID, boardID, listID, position)
		if err != nil {
			respondBasedOnError(c, err)
			return
		}
	}

	c.Status(http.StatusNoContent)
}
