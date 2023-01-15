package repository

import (
	"context"
	"fmt"
	"time"

	"firebase.google.com/go/v4/db"
	"github.com/jordyf15/thullo-api/card"
	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/jordyf15/thullo-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type cardRepository struct {
	dbClient *db.Client
}

func NewCardRepository(dbClient *db.Client) card.Repository {
	return &cardRepository{dbClient: dbClient}
}

func (repo *cardRepository) Create(card *models.Card) error {
	card.ID = primitive.NewObjectID()
	card.CreatedAt = time.Now()
	card.UpdatedAt = card.CreatedAt

	ctx := context.Background()
	ref := repo.dbClient.NewRef(fmt.Sprintf("cards/%s", card.ID.Hex()))

	return ref.Set(ctx, card)
}

func (repo *cardRepository) GetListCards(listID primitive.ObjectID) ([]*models.Card, error) {
	ctx := context.Background()
	ref := repo.dbClient.NewRef("cards").OrderByChild("list_id").EqualTo(listID.Hex())

	cardsMap := make(map[string]*models.Card)

	err := ref.Get(ctx, &cardsMap)
	if err != nil {
		return nil, err
	}

	cards := []*models.Card{}

	for _, card := range cardsMap {
		cards = append(cards, card)
	}

	return cards, nil
}

func (repo *cardRepository) GetCardByID(cardID primitive.ObjectID) (*models.Card, error) {
	ctx := context.Background()
	ref := repo.dbClient.NewRef(fmt.Sprintf("cards/%s", cardID.Hex()))

	card := &models.Card{}

	err := ref.Get(ctx, &card)
	if err != nil {
		return nil, err
	}

	if card == nil {
		return nil, custom_errors.ErrRecordNotFound
	}

	return card, nil
}
