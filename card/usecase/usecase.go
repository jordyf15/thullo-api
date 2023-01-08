package usecase

import (
	"github.com/jordyf15/thullo-api/card"
	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/jordyf15/thullo-api/list"
	"github.com/jordyf15/thullo-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type cardUsecase struct {
	listRepo list.Repository
	cardRepo card.Repository
}

func NewCardUsecase(listRepo list.Repository, cardRepo card.Repository) card.Usecase {
	return &cardUsecase{listRepo: listRepo, cardRepo: cardRepo}
}

func (usecase *cardUsecase) Create(listID primitive.ObjectID, title string) error {
	if title == "" {
		return custom_errors.ErrCardTitleEmpty
	}

	list, err := usecase.listRepo.GetListByID(listID)
	if err != nil {
		return err
	}

	cards, err := usecase.cardRepo.GetListCards(list.ID)
	if err != nil {
		return err
	}

	card := &models.Card{
		Title:    title,
		ListID:   list.ID,
		Position: len(cards),
	}

	err = usecase.cardRepo.Create(card)
	if err != nil {
		return err
	}

	return nil
}
