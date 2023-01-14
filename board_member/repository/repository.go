package repository

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/db"
	"github.com/jordyf15/thullo-api/board_member"
	"github.com/jordyf15/thullo-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type boardMemberRepository struct {
	dbClient *db.Client
}

func NewBoardMemberRepository(dbClient *db.Client) board_member.Repository {
	return &boardMemberRepository{dbClient: dbClient}
}

func (repo *boardMemberRepository) Create(boardMember *models.BoardMember) error {
	boardMember.ID = primitive.NewObjectID()

	ctx := context.Background()
	ref := repo.dbClient.NewRef(fmt.Sprintf("board_members/%s", boardMember.ID.Hex()))

	return ref.Set(ctx, boardMember)
}

func (repo *boardMemberRepository) GetBoardMembers(boardID primitive.ObjectID) ([]*models.BoardMember, error) {
	ctx := context.Background()
	ref := repo.dbClient.NewRef("board_members").OrderByChild("board_id").EqualTo(boardID.Hex())

	boardMembersMap := make(map[string]*models.BoardMember)

	err := ref.Get(ctx, &boardMembersMap)
	if err != nil {
		return nil, err
	}

	boardMembers := []*models.BoardMember{}

	for _, boardMember := range boardMembersMap {
		boardMembers = append(boardMembers, boardMember)
	}

	return boardMembers, nil
}

func (repo *boardMemberRepository) UpdateBoardMemberRole(ID primitive.ObjectID, role models.MemberRole) error {
	ctx := context.Background()
	ref := repo.dbClient.NewRef(fmt.Sprintf("board_members/%s/role", ID.Hex()))

	return ref.Set(ctx, role)
}
