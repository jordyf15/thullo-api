package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/jordyf15/thullo-api/models"
	"github.com/jordyf15/thullo-api/user"
	"github.com/jordyf15/thullo-api/user/repository"
	"github.com/jordyf15/thullo-api/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(userRepositorySuite))
}

type userRepositorySuite struct {
	suite.Suite
	collection *mongo.Collection
	repository user.Repository
}

var (
	userID1 = primitive.NewObjectID()

	user1 = &models.User{
		ID:                userID1,
		Email:             "user1@gmail.com",
		EncryptedPassword: "hashedPassword",
		Username:          "user1",
		Name:              "user1",
		Bio:               "first user",
		UpdatedAt:         time.Now(),
		CreatedAt:         time.Now(),
	}
)

func (s *userRepositorySuite) SetupSuite() {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	db := client.Database("thullo-test")
	s.collection = db.Collection("users")
	s.repository = repository.NewUserRepository(db)
}

func (s *userRepositorySuite) SetupTest() {
	s.collection.InsertOne(context.TODO(), utils.ToBSON(user1))
}

func (s *userRepositorySuite) AfterTest(suiteName, testName string) {
	s.collection.DeleteMany(context.TODO(), bson.M{})
}

func (s *userRepositorySuite) TestCreate() {
	user := &models.User{
		Email:             "newuser@gmail.com",
		EncryptedPassword: "hashedPassword",
		Username:          "newuser",
		Name:              "new user",
		Bio:               "a new user",
	}

	err := s.repository.Create(user)
	assert.NoError(s.T(), err)

	filter := bson.M{"username": "newuser"}
	var foundUser models.User
	err = s.collection.FindOne(context.TODO(), filter).Decode(&foundUser)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user.Email, foundUser.Email)
	assert.Equal(s.T(), user.EncryptedPassword, foundUser.EncryptedPassword)
	assert.Equal(s.T(), user.Username, foundUser.Username)
	assert.Equal(s.T(), user.Name, foundUser.Name)
	assert.Equal(s.T(), user.Bio, foundUser.Bio)
}

func (s *userRepositorySuite) TestFieldExistsTrue() {
	isExist, err := s.repository.FieldExists("username", "user1")
	assert.NoError(s.T(), err)
	assert.True(s.T(), isExist)
}

func (s *userRepositorySuite) TestFieldExistsFalse() {
	isExist, err := s.repository.FieldExists("username", "user15")
	assert.NoError(s.T(), err)
	assert.False(s.T(), isExist)
}

func (s *userRepositorySuite) TestGetByEmail() {
	user, err := s.repository.GetByEmail("user1@gmail.com")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user1.ID, user.ID)
	assert.Equal(s.T(), user1.Email, user.Email)
	assert.Equal(s.T(), user1.Username, user.Username)
	assert.Equal(s.T(), user1.Name, user.Name)
	assert.Equal(s.T(), user1.Bio, user.Bio)
}
