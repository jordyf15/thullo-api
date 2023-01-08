package models

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Card struct {
	ID          primitive.ObjectID `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	ListID      primitive.ObjectID `json:"list_id"`
	Cover       *BoardCover        `json:"cover"`
	Position    int                `json:"position"`
	// Attachments?
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (card *Card) MarshalJSON() ([]byte, error) {
	type Alias Card
	newStruct := &struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias: (*Alias)(card),
	}

	newStruct.CreatedAt = card.CreatedAt.Format("2006-01-02T15:04:05-0700")
	newStruct.UpdatedAt = card.UpdatedAt.Format("2006-01-02T15:04:05-0700")

	return json.Marshal(newStruct)
}

func (card *Card) UnmarshalJSON(data []byte) error {
	type Alias Card
	alias := &struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{Alias: (*Alias)(card)}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	card.CreatedAt, err = time.Parse("2006-01-02T15:04:05-0700", alias.CreatedAt)
	if err != nil {
		return err
	}

	card.UpdatedAt, err = time.Parse("2006-01-02T15:04:05-0700", alias.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
