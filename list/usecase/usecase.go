package usecase

import (
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
