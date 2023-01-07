package usecase

import (
	"time"

	"github.com/jordyf15/thullo-api/board"
	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/jordyf15/thullo-api/list"
	"github.com/jordyf15/thullo-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type listUsecase struct {
	listRepo  list.Repository
	boardRepo board.Repository
}

func NewListUsecase(listRepo list.Repository, boardRepo board.Repository) list.Usecase {
	return &listUsecase{listRepo: listRepo, boardRepo: boardRepo}
}

func (usecase *listUsecase) Create(boardID primitive.ObjectID, title string) error {
	if title == "" {
		return custom_errors.ErrListTitleEmpty
	}

	board, err := usecase.boardRepo.GetBoardByID(boardID)
	if err != nil {
		return err
	}

	lists, err := usecase.listRepo.GetBoardLists(board.ID)
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

func (usecase *listUsecase) UpdateTitle(listID primitive.ObjectID, title string) error {
	if title == "" {
		return custom_errors.ErrListTitleEmpty
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

func (usecase *listUsecase) UpdatePosition(boardID, listID primitive.ObjectID, newPosition int) error {
	updatedList, err := usecase.listRepo.GetListByID(listID)
	if err != nil {
		return err
	}
	prevPosition := updatedList.Position

	lists, err := usecase.listRepo.GetBoardLists(boardID)
	if err != nil {
		return err
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
