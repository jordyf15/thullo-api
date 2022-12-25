package repository

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/db"
	"github.com/jordyf15/thullo-api/models"
	"github.com/jordyf15/thullo-api/user_boards"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userBoardsRepository struct {
	dbClient *db.Client
}

func NewUserBoardsRepository(dbClient *db.Client) user_boards.Repository {
	return &userBoardsRepository{dbClient: dbClient}
}

func (repo *userBoardsRepository) Create(userID primitive.ObjectID) error {
	ctx := context.Background()
	ref := repo.dbClient.NewRef(fmt.Sprintf("user_boards/%s", userID.Hex()))

	userBoards := &models.UserBoards{
		UserID:   userID,
		BoardIDs: []primitive.ObjectID{},
	}

	err := ref.Set(ctx, userBoards)
	if err != nil {
		return err
	}

	return nil
}

func (repo *userBoardsRepository) AddBoard(userID, boardID primitive.ObjectID) error {
	ctx := context.Background()
	ref := repo.dbClient.NewRef(fmt.Sprintf("user_boards/%s/board_ids/%s", userID.Hex(), boardID.Hex()))

	err := ref.Set(ctx, true)
	if err != nil {
		return err
	}

	return nil
}
