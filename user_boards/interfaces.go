package user_boards

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	AddBoard(userID, boardID primitive.ObjectID) error
	Create(userID primitive.ObjectID) error
}
