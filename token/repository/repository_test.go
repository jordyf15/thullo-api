package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jordyf15/thullo-api/models"
	"github.com/jordyf15/thullo-api/token"
	"github.com/jordyf15/thullo-api/token/repository"
	"github.com/jordyf15/thullo-api/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestTokenRepository(t *testing.T) {
	suite.Run(t, new(tokenRepositorySuite))
}

type tokenRepositorySuite struct {
	suite.Suite
	collection *mongo.Collection
	repository token.Repository
	redis      *redis.Client
}

var (
	tokenID1            = primitive.NewObjectID()
	userID1             = primitive.NewObjectID()
	refreshTokenID1     = "222f1f77bcf86cd799439011"
	prevRefreshTokenID1 = "333f1f77bcf86cd799439011"
	tokenID2            = primitive.NewObjectID()
	userID2             = primitive.NewObjectID()
	refreshTokenID2     = "777f1f77bcf86cd799439011"
	prevRefreshTokenID2 = "888f1f77bcf86cd799439011"

	tokenSet1 = models.TokenSet{
		ID:                 tokenID1,
		UserID:             userID1,
		RefreshTokenID:     refreshTokenID1,
		PrevRefreshTokenID: &prevRefreshTokenID1,
		UpdatedAt:          time.Now(),
	}
	tokenSet2 = models.TokenSet{
		ID:                 tokenID2,
		UserID:             userID2,
		RefreshTokenID:     refreshTokenID2,
		PrevRefreshTokenID: &prevRefreshTokenID2,
		UpdatedAt:          time.Now(),
	}

	userID         = primitive.NewObjectID()
	refreshTokenID = "432f1f77bcf86cd799439011"
	accessToken    = models.AccessToken{
		UserID:         userID,
		RefreshTokenID: refreshTokenID,
		Token: models.Token{StandardClaims: jwt.StandardClaims{
			Id:        "54321",
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		}},
	}
)

func (s *tokenRepositorySuite) SetupSuite() {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	_miniredis, err := miniredis.Run()
	if err != nil {
		s.T().Fatalf("An error occured: %s", err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: _miniredis.Addr(),
	})
	db := client.Database("thullo-test")
	s.collection = db.Collection("token_sets")
	s.repository = repository.NewTokenRepository(db, redisClient)
	s.redis = redisClient
}

func (s *tokenRepositorySuite) SetupTest() {
	s.collection.InsertOne(context.TODO(), utils.ToBSON(tokenSet1))
	s.collection.InsertOne(context.TODO(), utils.ToBSON(tokenSet2))

	s.redis.ZAdd(context.TODO(), repository.RedisKeyFreshAccessTokens,
		redis.Z{Score: float64(accessToken.ExpiresAt), Member: accessToken.Id})
}

func (s *tokenRepositorySuite) AfterTest(suiteName, testName string) {
	s.collection.Drop(context.TODO())
	s.redis.FlushAll(context.TODO())
}

func (s *tokenRepositorySuite) TestCreate() {
	userId := primitive.NewObjectID()
	refreshTokenId := primitive.NewObjectID().Hex()
	prevRefreshTokenId := primitive.NewObjectID().Hex()
	tokenSet := models.TokenSet{
		UserID:             userId,
		RefreshTokenID:     refreshTokenId,
		PrevRefreshTokenID: &prevRefreshTokenId,
	}

	err := s.repository.Create(&tokenSet)

	assert.NoError(s.T(), err)

	foundTokenSet := &models.TokenSet{}
	filter := bson.D{
		{Key: "user_id", Value: userId},
	}
	s.collection.FindOne(context.TODO(), filter).Decode(foundTokenSet)

	assert.Equal(s.T(), userId.Hex(), foundTokenSet.UserID.Hex())
	assert.Equal(s.T(), refreshTokenId, foundTokenSet.RefreshTokenID)
	assert.Equal(s.T(), &prevRefreshTokenId, foundTokenSet.PrevRefreshTokenID)
}

func (s *tokenRepositorySuite) TestUpdate() {
	updatedRefreshTokenId1 := "123f1f77bcf86cd799439011"
	updatedPrevRefreshTokenId1 := "321f1f77bcf86cd799439011"
	updatedTokenSet1 := models.TokenSet{
		ID:                 tokenID1,
		UserID:             userID1,
		RefreshTokenID:     updatedRefreshTokenId1,
		PrevRefreshTokenID: &updatedPrevRefreshTokenId1,
		UpdatedAt:          time.Now(),
	}

	err := s.repository.Update(&updatedTokenSet1)
	assert.NoError(s.T(), err)

	foundTokenSet := &models.TokenSet{}
	filter := bson.D{
		{Key: "_id", Value: tokenID1},
	}

	s.collection.FindOne(context.TODO(), filter).Decode(foundTokenSet)

	assert.Equal(s.T(), tokenID1.Hex(), foundTokenSet.ID.Hex())
	assert.Equal(s.T(), userID1.Hex(), foundTokenSet.UserID.Hex())
	assert.Equal(s.T(), updatedRefreshTokenId1, foundTokenSet.RefreshTokenID)
	assert.Equal(s.T(), &updatedPrevRefreshTokenId1, foundTokenSet.PrevRefreshTokenID)
}

func (s *tokenRepositorySuite) TestDeleteRefreshTokenIdNotEmpty() {
	deletedTokenSet := models.TokenSet{
		ID:                 tokenID1,
		UserID:             userID1,
		RefreshTokenID:     refreshTokenID1,
		PrevRefreshTokenID: &prevRefreshTokenID1,
		UpdatedAt:          time.Now(),
	}

	err := s.repository.Delete(&deletedTokenSet)
	assert.NoError(s.T(), err)

	filter := bson.D{
		{Key: "_id", Value: tokenID1},
	}
	foundTokenSet := &models.TokenSet{}
	err = s.collection.FindOne(context.TODO(), filter).Decode(foundTokenSet)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), mongo.ErrNoDocuments, err)
}

func (s *tokenRepositorySuite) TestDeleteRefreshTokenIdEmpty() {
	prevRefreshTokenId := ""
	deletedTokenSet := models.TokenSet{
		ID:                 tokenID1,
		UserID:             userID1,
		RefreshTokenID:     "",
		PrevRefreshTokenID: &prevRefreshTokenId,
		UpdatedAt:          time.Now(),
	}

	err := s.repository.Delete(&deletedTokenSet)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), "Refresh token ID is empty", err.Error())
}

func (s *tokenRepositorySuite) TestUpdates() {
	updatedTokenSet := models.TokenSet{
		ID:                 tokenID1,
		UserID:             userID1,
		RefreshTokenID:     refreshTokenID1,
		PrevRefreshTokenID: &prevRefreshTokenID1,
		UpdatedAt:          time.Now(),
	}

	changes := map[string]interface{}{"prt_id": nil}

	err := s.repository.Updates(&updatedTokenSet, changes)
	assert.NoError(s.T(), err)

	foundTokenSet := &models.TokenSet{}
	filter := bson.D{
		{Key: "_id", Value: tokenID1},
	}
	s.collection.FindOne(context.TODO(), filter).Decode(foundTokenSet)
	assert.Equal(s.T(), tokenID1.Hex(), foundTokenSet.ID.Hex())
	assert.Equal(s.T(), userID1.Hex(), foundTokenSet.UserID.Hex())
	assert.Equal(s.T(), refreshTokenID1, foundTokenSet.RefreshTokenID)
	assert.Nil(s.T(), foundTokenSet.PrevRefreshTokenID)
}

func (s *tokenRepositorySuite) TestGetTokenSetIncludeParent() {
	foundTokenSet, err := s.repository.GetTokenSet(userID1, prevRefreshTokenID1, true)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), tokenID1.Hex(), foundTokenSet.ID.Hex())
	assert.Equal(s.T(), userID1.Hex(), foundTokenSet.UserID.Hex())
	assert.Equal(s.T(), refreshTokenID1, foundTokenSet.RefreshTokenID)
	assert.Equal(s.T(), prevRefreshTokenID1, *foundTokenSet.PrevRefreshTokenID)
}

func (s *tokenRepositorySuite) TestGetTokenSetNotIncludeParent() {
	foundTokenSet, err := s.repository.GetTokenSet(userID1, refreshTokenID1, false)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), tokenID1.Hex(), foundTokenSet.ID.Hex())
	assert.Equal(s.T(), userID1.Hex(), foundTokenSet.UserID.Hex())
	assert.Equal(s.T(), refreshTokenID1, foundTokenSet.RefreshTokenID)
	assert.Equal(s.T(), prevRefreshTokenID1, *foundTokenSet.PrevRefreshTokenID)
}

func (s *tokenRepositorySuite) TestGetTokenSetNotFound() {
	_, err := s.repository.GetTokenSet(userID, refreshTokenID1, false)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), mongo.ErrNoDocuments.Error(), err.Error())

	_, err = s.repository.GetTokenSet(userID, prevRefreshTokenID1, true)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), mongo.ErrNoDocuments.Error(), err.Error())
}

func (s *tokenRepositorySuite) TestSave() {
	s.redis.FlushAll(context.TODO())

	accessToken := models.AccessToken{
		UserID:         primitive.NewObjectID(),
		RefreshTokenID: primitive.NewObjectID().Hex(),
		Token: models.Token{StandardClaims: jwt.StandardClaims{
			Id:        "12345",
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		}},
	}

	err := s.repository.Save(&accessToken)
	assert.NoError(s.T(), err)

	accessTokenJson := s.redis.ZRange(context.TODO(), repository.RedisKeyFreshAccessTokens, 0, -1).Val()
	assert.Equal(s.T(), "12345", accessTokenJson[0])
}

func (s *tokenRepositorySuite) TestExistsTrue() {
	accessToken := models.AccessToken{
		UserID:         primitive.NewObjectID(),
		RefreshTokenID: primitive.NewObjectID().Hex(),
		Token: models.Token{StandardClaims: jwt.StandardClaims{
			Id:        "54321",
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		}},
	}

	isExist := s.repository.Exists(&accessToken)
	assert.True(s.T(), isExist)
}

func (s *tokenRepositorySuite) TestExistsFalse() {
	accessToken := models.AccessToken{
		UserID:         primitive.NewObjectID(),
		RefreshTokenID: primitive.NewObjectID().Hex(),
		Token: models.Token{StandardClaims: jwt.StandardClaims{
			Id:        "00000",
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		}},
	}

	isExist := s.repository.Exists(&accessToken)
	assert.False(s.T(), isExist)
}

func (s *tokenRepositorySuite) TestRemove() {
	accessToken := models.AccessToken{
		UserID:         userID,
		RefreshTokenID: refreshTokenID,
		Token: models.Token{StandardClaims: jwt.StandardClaims{
			Id:        "54321",
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		}},
	}

	err := s.repository.Remove(&accessToken)
	assert.NoError(s.T(), err)

	accessTokenJson := s.redis.ZRange(context.TODO(), repository.RedisKeyFreshAccessTokens, 0, -1).Val()
	assert.Len(s.T(), accessTokenJson, 0)
}
