package usecase

import (
	"github.com/jordyf15/thullo-api/board"
	"github.com/jordyf15/thullo-api/board_member"
	"github.com/jordyf15/thullo-api/card"
	"github.com/jordyf15/thullo-api/comment"
	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/jordyf15/thullo-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type commentUsecase struct {
	commentRepo     comment.Repository
	boardMemberRepo board_member.Repository
	boardRepo       board.Repository
	cardRepo        card.Repository
}

func NewCommentUsecase(boardMemberRepo board_member.Repository, cardRepo card.Repository, commentRepo comment.Repository, boardRepo board.Repository) comment.Usecase {
	return &commentUsecase{boardMemberRepo: boardMemberRepo, cardRepo: cardRepo, commentRepo: commentRepo, boardRepo: boardRepo}
}

func (usecase *commentUsecase) Create(requesterID, boardID, cardID primitive.ObjectID, comment string) error {
	if comment == "" {
		return custom_errors.ErrCommentEmpty
	}

	board, err := usecase.boardRepo.GetBoardByID(boardID)
	if err != nil {
		return err
	}

	if board.Visibility == models.BoardVisibilityPrivate {
		boardMembers, err := usecase.boardMemberRepo.GetBoardMembers(boardID)
		if err != nil {
			return err
		}

		isRequesterBoardMember := false
		for _, boardMember := range boardMembers {
			if boardMember.UserID == requesterID {
				isRequesterBoardMember = true
				break
			}
		}

		if !isRequesterBoardMember {
			return custom_errors.ErrNotAuthorized
		}
	}

	_, err = usecase.cardRepo.GetCardByID(cardID)
	if err != nil {
		return err
	}

	_comment := &models.Comment{
		AuthorID: requesterID,
		CardID:   cardID,
		Comment:  comment,
	}

	err = usecase.commentRepo.Create(_comment)
	if err != nil {
		return err
	}

	return nil
}
