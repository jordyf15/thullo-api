package board_member

import (
	"github.com/jordyf15/thullo-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	Create(boardMember *models.BoardMember) error
	GetBoardMembers(boardID primitive.ObjectID) ([]*models.BoardMember, error)
	UpdateBoardMemberRole(ID primitive.ObjectID, role models.MemberRole) error
	DeleteBoardMemberByID(ID primitive.ObjectID) error
}
