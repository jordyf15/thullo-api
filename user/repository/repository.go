package repository

import (
	"context"
	"time"

	"github.com/jordyf15/thullo-api/models"
	"github.com/jordyf15/thullo-api/user"
	"github.com/jordyf15/thullo-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const contextTimeout = time.Second * 30

type userRepository struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Database) user.Repository {
	collection := db.Collection("users")
	return &userRepository{db: collection}
}

func (repo *userRepository) Create(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := repo.db.InsertOne(ctx, utils.ToBSON(user))

	return err
}

func (repo *userRepository) FieldExists(key string, value string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	filter := bson.D{
		{Key: key, Value: value},
	}

	count, err := repo.db.CountDocuments(ctx, filter)
	if err != nil {
		return false, nil
	}

	return count > 0, nil
}

func (repo *userRepository) GetByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	filter := bson.D{
		{Key: "email", Value: email},
	}

	foundUser := &models.User{}
	err := repo.db.FindOne(ctx, filter).Decode(foundUser)

	return foundUser, err

}

func (repo *userRepository) GetByID(userID primitive.ObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	filter := bson.D{
		{Key: "_id", Value: userID},
	}

	foundUser := &models.User{}
	err := repo.db.FindOne(ctx, filter).Decode(foundUser)

	return foundUser, err
}
