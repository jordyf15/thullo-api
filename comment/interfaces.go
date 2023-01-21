package comment

import (
	"github.com/jordyf15/thullo-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	Create(comment *models.Comment) error
	GetCommentByID(commentID primitive.ObjectID) (*models.Comment, error)
	Update(comment *models.Comment) error
	DeleteCommentByID(commentID primitive.ObjectID) error
}

type Usecase interface {
	Create(requesterID, boardID, listID, cardID primitive.ObjectID, comment string) error
	Update(requesterID, boardID, listID, cardID, commentID primitive.ObjectID, comment string) error
	Delete(requesterID, boardID, listID, cardID, commentID primitive.ObjectID) error
}
