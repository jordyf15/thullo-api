package main

import (
	"net/http"

	"github.com/jordyf15/thullo-api/controllers"
	"github.com/jordyf15/thullo-api/storage"
	tr "github.com/jordyf15/thullo-api/token/repository"
	tu "github.com/jordyf15/thullo-api/token/usecase"

	bmr "github.com/jordyf15/thullo-api/board_member/repository"
	cr "github.com/jordyf15/thullo-api/card/repository"
	cmr "github.com/jordyf15/thullo-api/comment/repository"
	ur "github.com/jordyf15/thullo-api/user/repository"

	br "github.com/jordyf15/thullo-api/board/repository"
	bu "github.com/jordyf15/thullo-api/board/usecase"
	cu "github.com/jordyf15/thullo-api/card/usecase"
	cmu "github.com/jordyf15/thullo-api/comment/usecase"
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
	commentRepo := cmr.NewCommentRepository(rtdbClient)

	tokenUsecase := tu.NewTokenUsecase(tokenRepo)
	userUsecase := uu.NewUserUsecase(userRepo, tokenRepo, oauthRepo, _storage)
	boardUsecase := bu.NewBoardUsecase(boardRepo, unsplashRepo, boardMemberRepo, userRepo, _storage)
	listUsecase := lu.NewListUsecase(listRepo, boardRepo, boardMemberRepo)
	cardUsecase := cu.NewCardUsecase(listRepo, cardRepo, boardMemberRepo)
	commentUsecase := cmu.NewCommentUsecase(boardMemberRepo, cardRepo, commentRepo, boardRepo, listRepo)

	tokenController := controllers.NewTokenController(tokenUsecase)
	userController := controllers.NewUserController(userUsecase)
	boardController := controllers.NewBoardController(boardUsecase)
	listController := controllers.NewListController(listUsecase)
	cardController := controllers.NewCardController(cardUsecase)
	commentController := controllers.NewCommentController(commentUsecase)

	router.GET("_health", health)

	router.POST("tokens/refresh", tokenController.RefreshAccessToken)
	router.POST("tokens/remove", tokenController.DeleteRefreshToken)

	router.POST("register", userController.Register)
	router.POST("login", userController.Login)
	router.POST("login/google", userController.LoginWithGoogle)

	router.POST("boards", boardController.Create)
	router.PATCH("boards/:board_id", boardController.Update)

	router.POST("boards/:board_id/members", boardController.AddMember)
	router.PATCH("boards/:board_id/members/:member_id", boardController.UpdateMemberRole)
	router.DELETE("boards/:board_id/members/:member_id", boardController.DeleteMember)

	router.POST("boards/:board_id/lists", listController.Create)
	router.PATCH("boards/:board_id/lists/:list_id", listController.Update)

	router.POST("boards/:board_id/lists/:list_id/cards/:card_id/comments", commentController.Create)
	router.PATCH("boards/:board_id/lists/:list_id/cards/:card_id/comments/:comment_id", commentController.Update)
	router.DELETE("boards/:board_id/lists/:list_id/cards/:card_id/comments/:comment_id", commentController.Delete)

	router.POST("boards/:board_id/lists/:list_id/cards", cardController.Create)
}
