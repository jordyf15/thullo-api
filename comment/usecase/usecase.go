package usecase

import (
	"github.com/jordyf15/thullo-api/board"
	"github.com/jordyf15/thullo-api/board_member"
	"github.com/jordyf15/thullo-api/card"
	"github.com/jordyf15/thullo-api/comment"
	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/jordyf15/thullo-api/list"
	"github.com/jordyf15/thullo-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type commentUsecase struct {
	commentRepo     comment.Repository
	boardMemberRepo board_member.Repository
	boardRepo       board.Repository
	cardRepo        card.Repository
	listRepo        list.Repository
}

func NewCommentUsecase(boardMemberRepo board_member.Repository, cardRepo card.Repository, commentRepo comment.Repository, boardRepo board.Repository, listRepo list.Repository) comment.Usecase {
	return &commentUsecase{boardMemberRepo: boardMemberRepo, cardRepo: cardRepo, commentRepo: commentRepo, boardRepo: boardRepo, listRepo: listRepo}
}

func (usecase *commentUsecase) Create(requesterID, boardID, listID, cardID primitive.ObjectID, comment string) error {
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

	list, err := usecase.listRepo.GetListByID(listID)
	if err != nil {
		return err
	}

	if list.BoardID != boardID {
		return custom_errors.ErrRecordNotFound
	}

	card, err := usecase.cardRepo.GetCardByID(cardID)
	if err != nil {
		return err
	}

	if card.ListID != listID {
		return custom_errors.ErrRecordNotFound
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

func (usecase *commentUsecase) Update(requesterID, boardID, listID, cardID, commentID primitive.ObjectID, comment string) error {
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

	list, err := usecase.listRepo.GetListByID(listID)
	if err != nil {
		return err
	}

	if list.BoardID != boardID {
		return custom_errors.ErrRecordNotFound
	}

	card, err := usecase.cardRepo.GetCardByID(cardID)
	if err != nil {
		return err
	}

	if card.ListID != listID {
		return custom_errors.ErrRecordNotFound
	}

	commentObj, err := usecase.commentRepo.GetCommentByID(commentID)
	if err != nil {
		return err
	}

	if commentObj.CardID != cardID {
		return custom_errors.ErrRecordNotFound
	}

	if commentObj.AuthorID != requesterID {
		return custom_errors.ErrNotAuthorized
	}

	commentObj.Comment = comment

	err = usecase.commentRepo.Update(commentObj)
	if err != nil {
		return err
	}

	return nil
}

func (usecase *commentUsecase) Delete(requesterID, boardID, listID, cardID, commentID primitive.ObjectID) error {
	board, err := usecase.boardRepo.GetBoardByID(boardID)
	if err != nil {
		return err
	}

	boardMembers, err := usecase.boardMemberRepo.GetBoardMembers(boardID)
	if err != nil {
		return err
	}

	var requesterBoardMember *models.BoardMember
	for _, boardMember := range boardMembers {
		if boardMember.UserID == requesterID {
			requesterBoardMember = boardMember
		}
	}

	if board.Visibility == models.BoardVisibilityPrivate {
		if requesterBoardMember == nil {
			return custom_errors.ErrNotAuthorized
		}
	}

	list, err := usecase.listRepo.GetListByID(listID)
	if err != nil {
		return err
	}

	if list.BoardID != boardID {
		return custom_errors.ErrRecordNotFound
	}

	card, err := usecase.cardRepo.GetCardByID(cardID)
	if err != nil {
		return err
	}

	if card.ListID != listID {
		return custom_errors.ErrRecordNotFound
	}

	comment, err := usecase.commentRepo.GetCommentByID(commentID)
	if err != nil {
		return err
	}

	if comment.CardID != cardID {
		return custom_errors.ErrRecordNotFound
	}

	if comment.AuthorID != requesterID && (requesterBoardMember == nil || requesterBoardMember.Role != models.MemberRoleAdmin) {
		return custom_errors.ErrNotAuthorized
	}

	err = usecase.commentRepo.DeleteCommentByID(commentID)
	if err != nil {
		return err
	}

	return nil
}
