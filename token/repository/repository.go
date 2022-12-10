package repository

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/jordyf15/thullo-api/models"
	"github.com/jordyf15/thullo-api/token"
	"github.com/jordyf15/thullo-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	RedisKeyFreshAccessTokens = "fresh-access-tokens"
	contextTimeout            = time.Second * 30
)

type tokenRepository struct {
	db    *mongo.Collection
	redis *redis.Client
}

func NewTokenRepository(db *mongo.Database, redis *redis.Client) token.Repository {
	collection := db.Collection("token_sets")

	return &tokenRepository{db: collection, redis: redis}
}

func (repo *tokenRepository) GetTokenSet(userID primitive.ObjectID, hashedRefreshTokenID string, includeParent bool) (*models.TokenSet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	tokenSet := &models.TokenSet{}
	filter := bson.D{}
	var err error
	if includeParent {
		filter = bson.D{
			{Key: "user_id", Value: userID},
			{Key: "$or", Value: []interface{}{
				bson.D{{Key: "rt_id", Value: hashedRefreshTokenID}},
				bson.D{{Key: "prt_id", Value: hashedRefreshTokenID}},
			},
			}}
		err = repo.db.FindOne(ctx, filter).Decode(tokenSet)
	} else {
		filter = bson.D{
			{Key: "user_id", Value: userID},
			{Key: "rt_id", Value: hashedRefreshTokenID},
		}
		err = repo.db.FindOne(ctx, filter).Decode(tokenSet)
	}

	return tokenSet, err
}

func (repo *tokenRepository) Save(accessToken *models.AccessToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	return repo.redis.ZAdd(ctx, RedisKeyFreshAccessTokens,
		redis.Z{Score: float64(accessToken.ExpiresAt), Member: accessToken.Id}).Err()
}

func (repo *tokenRepository) Exists(accessToken *models.AccessToken) bool {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	expiry := repo.redis.ZScore(ctx, RedisKeyFreshAccessTokens, accessToken.Id)
	return expiry != nil && expiry.Val() > 0
}

func (repo *tokenRepository) Remove(accessToken *models.AccessToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	return repo.redis.ZRem(ctx, RedisKeyFreshAccessTokens, accessToken.Id).Err()
}

func (repo *tokenRepository) Create(tokenSet *models.TokenSet) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	tokenSet.ID = primitive.NewObjectID()
	tokenSet.UpdatedAt = time.Now()

	_, err := repo.db.InsertOne(ctx, utils.ToBSON(tokenSet))

	return err
}

func (repo *tokenRepository) Update(tokenSet *models.TokenSet) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	updates := bson.D{
		{
			Key: "$set", Value: bson.D{
				{Key: "updated_at", Value: time.Now()},
				{Key: "rt_id", Value: tokenSet.RefreshTokenID},
				{Key: "prt_id", Value: tokenSet.PrevRefreshTokenID},
			},
		},
	}

	_, err := repo.db.UpdateByID(ctx, tokenSet.ID, updates)

	return err
}

func (repo *tokenRepository) Updates(tokenSet *models.TokenSet, changes map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	setValues := bson.D{{Key: "updated_at", Value: time.Now()}}
	for k, v := range changes {
		setValues = append(setValues, bson.E{Key: k, Value: v})
	}

	updates := bson.D{{Key: "$set", Value: setValues}}
	_, err := repo.db.UpdateByID(ctx, tokenSet.ID, updates)

	return err
}

func (repo *tokenRepository) Delete(tokenSet *models.TokenSet) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	if len(tokenSet.RefreshTokenID) > 0 {
		filter := bson.D{{Key: "user_id", Value: tokenSet.UserID}, {Key: "rt_id", Value: tokenSet.RefreshTokenID}}
		_, err := repo.db.DeleteMany(ctx, filter)
		return err
	}

	return errors.New("Refresh token ID is empty")
}
