package main

import (
	"github.com/jordyf15/thullo-api/controllers"
	tr "github.com/jordyf15/thullo-api/token/repository"
	tu "github.com/jordyf15/thullo-api/token/usecase"
)

func initializeRoutes() {
	tokenRepo := tr.NewTokenRepository(dbClient, redisClient)

	tokenUsecase := tu.NewTokenUsecase(tokenRepo)

	tokenController := controllers.NewTokenController(tokenUsecase)

	router.GET("_health", health)

	router.POST("tokens/refresh", tokenController.RefreshAccessToken)
	router.POST("tokens/remove", tokenController.DeleteRefreshToken)
}
