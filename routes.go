package main

import (
	"net/http"

	"github.com/jordyf15/thullo-api/controllers"
	"github.com/jordyf15/thullo-api/storage"
	tr "github.com/jordyf15/thullo-api/token/repository"
	tu "github.com/jordyf15/thullo-api/token/usecase"

	ur "github.com/jordyf15/thullo-api/user/repository"
	uu "github.com/jordyf15/thullo-api/user/usecase"

	or "github.com/jordyf15/thullo-api/oauth/repository"
)

func initializeRoutes() {
	_storage := storage.NewImgurStorage(&http.Client{})

	tokenRepo := tr.NewTokenRepository(dbClient, redisClient)
	userRepo := ur.NewUserRepository(dbClient)
	oauthRepo := or.NewOauthRepository(&http.Client{})

	tokenUsecase := tu.NewTokenUsecase(tokenRepo)
	userUsecase := uu.NewUserUsecase(userRepo, tokenRepo, oauthRepo, _storage)

	tokenController := controllers.NewTokenController(tokenUsecase)
	userController := controllers.NewUserController(userUsecase)

	router.GET("_health", health)

	router.POST("tokens/refresh", tokenController.RefreshAccessToken)
	router.POST("tokens/remove", tokenController.DeleteRefreshToken)

	router.POST("register", userController.Register)
	router.POST("login", userController.Login)
	router.POST("login/google", userController.LoginWithGoogle)
}
