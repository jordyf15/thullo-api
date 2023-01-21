package repository

import (
	"context"
	"fmt"
	"time"

	"firebase.google.com/go/v4/db"
	"github.com/jordyf15/thullo-api/comment"
	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/jordyf15/thullo-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type commentRepository struct {
	dbClient *db.Client
}

func NewCommentRepository(dbClient *db.Client) comment.Repository {
	return &commentRepository{dbClient: dbClient}
}

func (repo *commentRepository) Create(comment *models.Comment) error {
	comment.ID = primitive.NewObjectID()
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = comment.CreatedAt

	ctx := context.Background()
	ref := repo.dbClient.NewRef(fmt.Sprintf("comments/%s", comment.ID.Hex()))

	return ref.Set(ctx, comment)
}

func (repo *commentRepository) GetCommentByID(commentID primitive.ObjectID) (*models.Comment, error) {
	ctx := context.Background()
	ref := repo.dbClient.NewRef(fmt.Sprintf("comments/%s", commentID.Hex()))

	comment := &models.Comment{}

	err := ref.Get(ctx, &comment)
	if err != nil {
		return nil, err
	}

	if comment == nil {
		return nil, custom_errors.ErrRecordNotFound
	}

	return comment, nil
}

func (repo *commentRepository) Update(comment *models.Comment) error {
	comment.UpdatedAt = time.Now()

	ctx := context.Background()
	ref := repo.dbClient.NewRef(fmt.Sprintf("comments/%s", comment.ID.Hex()))

	return ref.Set(ctx, comment)
}

func (repo *commentRepository) DeleteCommentByID(commentID primitive.ObjectID) error {
	ctx := context.Background()
	ref := repo.dbClient.NewRef(fmt.Sprintf("comments/%s", commentID.Hex()))

	return ref.Delete(ctx)
}
