package main

import (
	"net/http"

	"github.com/jordyf15/thullo-api/controllers"
	"github.com/jordyf15/thullo-api/storage"
	tr "github.com/jordyf15/thullo-api/token/repository"
	tu "github.com/jordyf15/thullo-api/token/usecase"

	ur "github.com/jordyf15/thullo-api/user/repository"

	br "github.com/jordyf15/thullo-api/board/repository"
	bu "github.com/jordyf15/thullo-api/board/usecase"
	unr "github.com/jordyf15/thullo-api/unsplash/repository"
	uu "github.com/jordyf15/thullo-api/user/usecase"

	or "github.com/jordyf15/thullo-api/oauth/repository"
)

func initializeRoutes() {
	_storage := storage.NewImgurStorage(&http.Client{})

	tokenRepo := tr.NewTokenRepository(dbClient, redisClient)
	userRepo := ur.NewUserRepository(dbClient)
	oauthRepo := or.NewOauthRepository(&http.Client{})
	boardRepo := br.NewBoardRepository(rtdbClient)
	unsplashRepo := unr.NewUnsplashRepository(&http.Client{})

	tokenUsecase := tu.NewTokenUsecase(tokenRepo)
	userUsecase := uu.NewUserUsecase(userRepo, tokenRepo, oauthRepo, _storage)
	boardUsecase := bu.NewBoardUsecase(boardRepo, unsplashRepo, _storage)

	tokenController := controllers.NewTokenController(tokenUsecase)
	userController := controllers.NewUserController(userUsecase)
	boardController := controllers.NewBoardController(boardUsecase)

	router.GET("_health", health)

	router.POST("tokens/refresh", tokenController.RefreshAccessToken)
	router.POST("tokens/remove", tokenController.DeleteRefreshToken)

	router.POST("register", userController.Register)
	router.POST("login", userController.Login)
	router.POST("login/google", userController.LoginWithGoogle)

	router.POST("boards", boardController.Create)
}
