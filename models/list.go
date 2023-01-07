package models

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type List struct {
	ID        primitive.ObjectID `json:"id"`
	Title     string             `json:"title"`
	BoardID   primitive.ObjectID `json:"board_id"`
	Position  int                `json:"position"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

func (list *List) MarshalJSON() ([]byte, error) {
	type Alias List
	newStruct := &struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias: (*Alias)(list),
	}

	newStruct.CreatedAt = list.CreatedAt.Format("2006-01-02T15:04:05-0700")
	newStruct.UpdatedAt = list.UpdatedAt.Format("2006-01-02T15:04:05-0700")

	return json.Marshal(newStruct)
}

func (list *List) UnmarshalJSON(data []byte) error {
	type Alias List
	alias := &struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{Alias: (*Alias)(list)}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	list.CreatedAt, err = time.Parse("2006-01-02T15:04:05-0700", alias.CreatedAt)
	if err != nil {
		return err
	}

	list.UpdatedAt, err = time.Parse("2006-01-02T15:04:05-0700", alias.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
