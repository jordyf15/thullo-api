package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/thullo-api/models"
	"github.com/jordyf15/thullo-api/user"
)

type UserController interface {
	Register(c *gin.Context)
	LoginWithGoogle(c *gin.Context)
	Login(c *gin.Context)
}

type userController struct {
	userUsecase user.Usecase
}

func NewUserController(userUsecase user.Usecase) UserController {
	return &userController{userUsecase: userUsecase}
}

func (controller *userController) Register(c *gin.Context) {
	user := &models.User{}
	user.Name = c.PostForm("name")
	user.Username = c.PostForm("username")
	user.Email = c.PostForm("email")
	user.Password = c.PostForm("password")

	resp, err := controller.userUsecase.Create(user, nil)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (controller *userController) LoginWithGoogle(c *gin.Context) {
	loginResponse, err := controller.userUsecase.LoginWithGoogle(c.PostForm("token"))
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.JSON(http.StatusOK, loginResponse)
}

func (controller *userController) Login(c *gin.Context) {
	response, err := controller.userUsecase.Login(c.PostForm("email"), c.PostForm("password"))
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
