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
	AddMember(c *gin.Context)
	UpdateMemberRole(c *gin.Context)
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
	title := strings.TrimSpace(c.PostForm("title"))
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

func (controller *boardController) AddMember(c *gin.Context) {
	requesterID := c.MustGet("current_user_id").(primitive.ObjectID)
	boardIDStr := c.Param("board_id")
	memberIDStr := c.PostForm("member_id")

	boardID, err := primitive.ObjectIDFromHex(boardIDStr)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	memberID, err := primitive.ObjectIDFromHex(memberIDStr)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	err = controller.usecase.AddMember(requesterID, boardID, memberID)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (controller *boardController) UpdateMemberRole(c *gin.Context) {
	requesterID := c.MustGet("current_user_id").(primitive.ObjectID)
	boardIDStr := c.Param("board_id")
	memberIDStr := c.Param("member_id")
	role := c.PostForm("role")

	boardID, err := primitive.ObjectIDFromHex(boardIDStr)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	memberID, err := primitive.ObjectIDFromHex(memberIDStr)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	err = controller.usecase.UpdateMemberRole(requesterID, boardID, memberID, role)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
