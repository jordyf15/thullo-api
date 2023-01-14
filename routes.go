package main

import (
	"net/http"

	"github.com/jordyf15/thullo-api/controllers"
	"github.com/jordyf15/thullo-api/storage"
	tr "github.com/jordyf15/thullo-api/token/repository"
	tu "github.com/jordyf15/thullo-api/token/usecase"

	bmr "github.com/jordyf15/thullo-api/board_member/repository"
	cr "github.com/jordyf15/thullo-api/card/repository"
	ur "github.com/jordyf15/thullo-api/user/repository"

	br "github.com/jordyf15/thullo-api/board/repository"
	bu "github.com/jordyf15/thullo-api/board/usecase"
	cu "github.com/jordyf15/thullo-api/card/usecase"
	lr "github.com/jordyf15/thullo-api/list/repository"
	lu "github.com/jordyf15/thullo-api/list/usecase"
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
	listRepo := lr.NewListRepository(rtdbClient)
	cardRepo := cr.NewCardRepository(rtdbClient)
	boardMemberRepo := bmr.NewBoardMemberRepository(rtdbClient)

	tokenUsecase := tu.NewTokenUsecase(tokenRepo)
	userUsecase := uu.NewUserUsecase(userRepo, tokenRepo, oauthRepo, _storage)
	boardUsecase := bu.NewBoardUsecase(boardRepo, unsplashRepo, boardMemberRepo, userRepo, _storage)
	listUsecase := lu.NewListUsecase(listRepo, boardRepo)
	cardUsecase := cu.NewCardUsecase(listRepo, cardRepo)

	tokenController := controllers.NewTokenController(tokenUsecase)
	userController := controllers.NewUserController(userUsecase)
	boardController := controllers.NewBoardController(boardUsecase)
	listController := controllers.NewListController(listUsecase)
	cardController := controllers.NewCardController(cardUsecase)

	router.GET("_health", health)

	router.POST("tokens/refresh", tokenController.RefreshAccessToken)
	router.POST("tokens/remove", tokenController.DeleteRefreshToken)

	router.POST("register", userController.Register)
	router.POST("login", userController.Login)
	router.POST("login/google", userController.LoginWithGoogle)

	router.POST("boards", boardController.Create)
	router.POST("boards/:board_id/members", boardController.AddMember)
	router.PATCH("boards/:board_id/members/:member_id", boardController.UpdateMemberRole)

	router.POST("boards/:board_id/lists", listController.Create)
	router.PATCH("boards/:board_id/lists/:list_id", listController.Update)

	router.POST("boards/:board_id/lists/:list_id/cards", cardController.Create)
}
