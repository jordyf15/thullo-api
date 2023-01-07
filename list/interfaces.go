package list

import (
	"github.com/jordyf15/thullo-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	Create(list *models.List) error
	GetBoardLists(boardID primitive.ObjectID) ([]*models.List, error)
}

type Usecase interface {
	Create(boardID primitive.ObjectID, title string) error
}
