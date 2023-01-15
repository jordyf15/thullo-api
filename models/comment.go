package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `json:"id"`
	AuthorID  primitive.ObjectID `json:"author_id"`
	CardID    primitive.ObjectID `json:"card_id"`
	Comment   string             `json:"comment"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
