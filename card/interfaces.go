package card

import (
	"github.com/jordyf15/thullo-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	Create(card *models.Card) error
	GetListCards(listID primitive.ObjectID) ([]*models.Card, error)
}

type Usecase interface {
	Create(listID primitive.ObjectID, title string) error
}
