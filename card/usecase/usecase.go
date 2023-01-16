package usecase

import (
	"github.com/jordyf15/thullo-api/board_member"
	"github.com/jordyf15/thullo-api/card"
	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/jordyf15/thullo-api/list"
	"github.com/jordyf15/thullo-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type cardUsecase struct {
	listRepo        list.Repository
	cardRepo        card.Repository
	boardMemberRepo board_member.Repository
}

func NewCardUsecase(listRepo list.Repository, cardRepo card.Repository, boardMemberRepo board_member.Repository) card.Usecase {
	return &cardUsecase{listRepo: listRepo, cardRepo: cardRepo, boardMemberRepo: boardMemberRepo}
}

func (usecase *cardUsecase) Create(requesterID, boardID, listID primitive.ObjectID, title string) error {
	if title == "" {
		return custom_errors.ErrCardTitleEmpty
	}

	err := usecase.checkIfRequesterIsMemberOfBoard(requesterID, boardID)
	if err != nil {
		return err
	}

	list, err := usecase.listRepo.GetListByID(listID)
	if err != nil {
		return err
	}

	// make sure the list actually belong to the board that the user have access to
	if list.BoardID != boardID {
		return custom_errors.ErrRecordNotFound
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

func (usecase *cardUsecase) checkIfRequesterIsMemberOfBoard(requesterID, boardID primitive.ObjectID) error {
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
