package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/thullo-api/board"
	"github.com/jordyf15/thullo-api/custom_errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoardController interface {
	Create(c *gin.Context)
}

type boardController struct {
	usecase board.Usecase
}

func NewBoardController(usecase board.Usecase) BoardController {
	return &boardController{usecase: usecase}
}

func (controller *boardController) Create(c *gin.Context) {
	userID := c.MustGet("current_user_id").(primitive.ObjectID)
	coverString := c.PostForm("cover")
	title := c.PostForm("title")
	visibility := c.PostForm("visibility")

	boardCover := map[string]interface{}{}
	var err error
	if len(coverString) > 0 {
		coverSlice := strings.Split(coverString, ":")
		if len(coverSlice) != 3 {
			respondBasedOnError(c, custom_errors.ErrMalformedCover)
			return
		}

		boardCover["source"] = coverSlice[0]
		boardCover["photo_id"] = coverSlice[1]
		boardCover["fp_y"], err = strconv.ParseFloat(coverSlice[2], 64)
		if err != nil {
			respondBasedOnError(c, custom_errors.ErrMalformedCover)
			return
		}
	} else {
		respondBasedOnError(c, custom_errors.ErrBoardCoverEmpty)
		return
	}

	err = controller.usecase.Create(userID, title, visibility, boardCover)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
