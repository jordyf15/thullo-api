package list

import (
	"github.com/jordyf15/thullo-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	Create(list *models.List) error
	GetBoardLists(boardID primitive.ObjectID) ([]*models.List, error)
	GetListByID(listID primitive.ObjectID) (*models.List, error)
	UpdateList(listID primitive.ObjectID, list *models.List) error
}

type Usecase interface {
	Create(requesterID, boardID primitive.ObjectID, title string) error
	UpdateTitle(requesterID, boardID, listID primitive.ObjectID, title string) error
	UpdatePosition(requesterID, boardID, listID primitive.ObjectID, newPosition int) error
}
