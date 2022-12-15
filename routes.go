package main

import (
	"net/http"

	"github.com/jordyf15/thullo-api/controllers"
	"github.com/jordyf15/thullo-api/storage"
	tr "github.com/jordyf15/thullo-api/token/repository"
	tu "github.com/jordyf15/thullo-api/token/usecase"

	ur "github.com/jordyf15/thullo-api/user/repository"
	uu "github.com/jordyf15/thullo-api/user/usecase"
)

func initializeRoutes() {
	_storage := storage.NewImgurStorage(&http.Client{})

	tokenRepo := tr.NewTokenRepository(dbClient, redisClient)
	userRepo := ur.NewUserRepository(dbClient)

	tokenUsecase := tu.NewTokenUsecase(tokenRepo)
	userUsecase := uu.NewUserUsecase(userRepo, tokenRepo, _storage)

	tokenController := controllers.NewTokenController(tokenUsecase)
	userController := controllers.NewUserController(userUsecase)

	router.GET("_health", health)

	router.POST("tokens/refresh", tokenController.RefreshAccessToken)
	router.POST("tokens/remove", tokenController.DeleteRefreshToken)

	router.POST("register", userController.Register)
}
