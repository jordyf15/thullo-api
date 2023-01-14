package board

import (
	"github.com/jordyf15/thullo-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	BoardCoverSources = map[string]bool{
		"unsplash": true,
	}
	BoardCoverSizes = []uint{1080, 450, 150}
)

type Repository interface {
	Create(board *models.Board) error
	GetBoardByID(boardID primitive.ObjectID) (*models.Board, error)
}

type Usecase interface {
	Create(userID primitive.ObjectID, title string, visibility string, boardCover map[string]interface{}) error
	AddMember(requesterID, boardID, memberID primitive.ObjectID) error
	UpdateMemberRole(requesterID, boardID, memberID primitive.ObjectID, role string) error
}
