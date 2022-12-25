package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserBoards struct {
	UserID   primitive.ObjectID   `json:"user_id"`
	BoardIDs []primitive.ObjectID `json:"board_ids"`
}
