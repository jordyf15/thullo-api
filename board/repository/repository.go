package repository

import (
	"context"
	"fmt"
	"time"

	"firebase.google.com/go/v4/db"
	"github.com/jordyf15/thullo-api/board"
	"github.com/jordyf15/thullo-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type boardRepository struct {
	dbClient *db.Client
}

func NewBoardRepository(dbClient *db.Client) board.Repository {
	return &boardRepository{dbClient: dbClient}
}

func (repo *boardRepository) Create(board *models.Board) error {
	board.ID = primitive.NewObjectID()
	board.CreatedAt = time.Now()
	board.UpdatedAt = board.CreatedAt

	ctx := context.Background()
	ref := repo.dbClient.NewRef(fmt.Sprintf("boards/%s", board.ID.Hex()))

	return ref.Set(ctx, board)
}
