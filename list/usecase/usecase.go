package usecase

import (
	"time"

	"github.com/jordyf15/thullo-api/board"
	"github.com/jordyf15/thullo-api/board_member"
	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/jordyf15/thullo-api/list"
	"github.com/jordyf15/thullo-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type listUsecase struct {
	listRepo        list.Repository
	boardRepo       board.Repository
	boardMemberRepo board_member.Repository
}

func NewListUsecase(listRepo list.Repository, boardRepo board.Repository, boardMemberRepo board_member.Repository) list.Usecase {
	return &listUsecase{listRepo: listRepo, boardRepo: boardRepo, boardMemberRepo: boardMemberRepo}
}

func (usecase *listUsecase) Create(requesterID, boardID primitive.ObjectID, title string) error {
	if title == "" {
		return custom_errors.ErrListTitleEmpty
	}

	err := usecase.checkIfRequesterIsMemberOfBoard(requesterID, boardID)
	if err != nil {
		return err
	}

	lists, err := usecase.listRepo.GetBoardLists(boardID)
	if err != nil {
		return err
	}

	list := &models.List{
		Title:    title,
		BoardID:  boardID,
		Position: len(lists),
	}

	err = usecase.listRepo.Create(list)
	if err != nil {
		return err
	}

	return nil
}

func (usecase *listUsecase) UpdateTitle(requesterID, boardID, listID primitive.ObjectID, title string) error {
	if title == "" {
		return custom_errors.ErrListTitleEmpty
	}

	err := usecase.checkIfRequesterIsMemberOfBoard(requesterID, boardID)
	if err != nil {
		return err
	}

	list, err := usecase.listRepo.GetListByID(listID)
	if err != nil {
		return err
	}

	list.Title = title
	list.UpdatedAt = time.Now()

	err = usecase.listRepo.UpdateList(listID, list)
	if err != nil {
		return err
	}

	return nil
}

func (usecase *listUsecase) UpdatePosition(requesterID, boardID, listID primitive.ObjectID, newPosition int) error {
	err := usecase.checkIfRequesterIsMemberOfBoard(requesterID, boardID)
	if err != nil {
		return err
	}

	updatedList, err := usecase.listRepo.GetListByID(listID)
	if err != nil {
		return err
	}
	prevPosition := updatedList.Position

	lists, err := usecase.listRepo.GetBoardLists(boardID)
	if err != nil {
		return err
	}

	if newPosition < 0 {
		return custom_errors.ErrListPositionTooLow
	}

	if newPosition >= len(lists) {
		return custom_errors.ErrListPositionTooHigh
	}

	for _, list := range lists {
		if updatedList.ID == list.ID {
			updatedList.Position = newPosition

			err = usecase.listRepo.UpdateList(updatedList.ID, updatedList)
			if err != nil {
				return err
			}
		} else {
			isChange := usecase.readjustOtherListPosition(list, prevPosition, newPosition)

			if isChange {
				err = usecase.listRepo.UpdateList(list.ID, list)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (usecase *listUsecase) checkIfRequesterIsMemberOfBoard(requesterID, boardID primitive.ObjectID) error {
	boardMembers, err := usecase.boardMemberRepo.GetBoardMembers(boardID)
	if err != nil {
		return err
	}

	// if there are no board members it means there are no board
	// a board will always atleast have 1 member
	if len(boardMembers) == 0 {
		return custom_errors.ErrRecordNotFound
	}

	var requesterBoardMember *models.BoardMember

	for _, boardMember := range boardMembers {
		if boardMember.UserID == requesterID {
			requesterBoardMember = boardMember
		}
	}

	if requesterBoardMember == nil {
		return custom_errors.ErrNotAuthorized
	}

	return nil
}

func (usecase *listUsecase) readjustOtherListPosition(otherList *models.List, prevPosition, newPosition int) bool {
	if newPosition > prevPosition {
		if otherList.Position > prevPosition && otherList.Position <= newPosition {
			otherList.Position -= 1
			return true
		}
	} else {
		if otherList.Position < prevPosition && otherList.Position >= newPosition {
			otherList.Position += 1
			return true
		}
	}

	return false
}
