package repository

import (
	"context"
	"fmt"
	"time"

	"firebase.google.com/go/v4/db"
	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/jordyf15/thullo-api/list"
	"github.com/jordyf15/thullo-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type listRepository struct {
	dbClient *db.Client
}

func NewListRepository(dbClient *db.Client) list.Repository {
	return &listRepository{dbClient: dbClient}
}

func (repo *listRepository) Create(list *models.List) error {
	list.ID = primitive.NewObjectID()
	list.CreatedAt = time.Now()
	list.UpdatedAt = list.CreatedAt

	ctx := context.Background()
	ref := repo.dbClient.NewRef(fmt.Sprintf("lists/%s", list.ID.Hex()))

	return ref.Set(ctx, list)
}

func (repo *listRepository) GetBoardLists(boardID primitive.ObjectID) ([]*models.List, error) {
	ctx := context.Background()
	ref := repo.dbClient.NewRef("lists").OrderByChild("board_id").EqualTo(boardID.Hex())

	listsMap := make(map[string]*models.List)

	err := ref.Get(ctx, &listsMap)

	if err != nil {
		return nil, err
	}
	lists := []*models.List{}

	for _, list := range listsMap {
		lists = append(lists, list)
	}

	return lists, nil
}

func (repo *listRepository) GetListByID(listID primitive.ObjectID) (*models.List, error) {
	ctx := context.Background()
	ref := repo.dbClient.NewRef(fmt.Sprintf("lists/%s", listID.Hex()))

	list := &models.List{}

	err := ref.Get(ctx, &list)
	if err != nil {
		return nil, err
	}

	if list == nil {
		return nil, custom_errors.ErrRecordNotFound
	}

	return list, nil
}

func (repo *listRepository) UpdateList(listID primitive.ObjectID, list *models.List) error {
	ctx := context.Background()
	ref := repo.dbClient.NewRef(fmt.Sprintf("lists/%s", listID.Hex()))

	return ref.Set(ctx, list)
}
