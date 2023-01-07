package models

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoardVisibility string

const (
	BoardVisibilityPublic  = "public"
	BoardVisibilityPrivate = "private"
)

type Board struct {
	ID          primitive.ObjectID `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Visibility  BoardVisibility    `json:"visibility"`
	OwnerID     primitive.ObjectID `json:"owner_id"`
	Cover       *BoardCover        `json:"cover"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type BoardCover struct {
	PhotoID     string  `json:"photo_id"`
	Source      string  `json:"source"`
	FocalPointY float64 `json:"fp_y"`
	Images      Images  `json:"images"`
}

func (board *Board) SetVisibility(visibility string) {
	switch visibility {
	case BoardVisibilityPrivate:
		board.Visibility = BoardVisibilityPrivate
	case BoardVisibilityPublic:
		board.Visibility = BoardVisibilityPublic
	}
}

func (board *Board) EmptyImageURLs() {
	for _, img := range board.Cover.Images {
		img.URL = ""
	}
}

func (board *Board) MarshalJSON() ([]byte, error) {
	type Alias Board
	newStruct := &struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias: (*Alias)(board),
	}

	newStruct.CreatedAt = board.CreatedAt.Format("2006-01-02T15:04:05-0700")
	newStruct.UpdatedAt = board.UpdatedAt.Format("2006-01-02T15:04:05-0700")

	return json.Marshal(newStruct)
}

func (board *Board) UnmarshalJSON(data []byte) error {
	type Alias Board
	alias := &struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{Alias: (*Alias)(board)}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	board.CreatedAt, err = time.Parse("2006-01-02T15:04:05-0700", alias.CreatedAt)
	if err != nil {
		return err
	}

	board.UpdatedAt, err = time.Parse("2006-01-02T15:04:05-0700", alias.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
