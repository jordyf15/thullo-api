package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/thullo-api/models"
	"github.com/jordyf15/thullo-api/user"
)

type UserController interface {
	Register(c *gin.Context)
}

type userController struct {
	userUsecase user.Usecase
}

const (
	maxPictureSize = 5 * 1024 * 1024
)

func NewUserController(userUsecase user.Usecase) UserController {
	return &userController{userUsecase: userUsecase}
}

func (controller *userController) Register(c *gin.Context) {
	user := &models.User{}
	user.Name = c.PostForm("name")
	user.Username = c.PostForm("username")
	user.Email = c.PostForm("email")
	user.Password = c.PostForm("password")

	resp, err := controller.userUsecase.Create(user)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}