package models

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	jwt.StandardClaims
}

type AccessToken struct {
	UserID         primitive.ObjectID `json:"uid"`
	RefreshTokenID string             `json:"rt_id"`
	Token
}

func (token *AccessToken) SetExpiration(expiryTime time.Time) *AccessToken {
	token.ExpiresAt = expiryTime.Unix()
	return token
}

func (token *AccessToken) ToJWT() *jwt.Token {
	return jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), token)
}

func (token *AccessToken) ToJWTString() string {
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), token)
	tokenString, _ := jwtToken.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
	return tokenString
}

type RefreshToken struct {
	UserID primitive.ObjectID `json:"uid"`
	Token
}

func (refreshToken *RefreshToken) ToJWT() *jwt.Token {
	return jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), refreshToken)
}

func (refreshToken *RefreshToken) ToJWTString() string {
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), refreshToken)
	tokenString, _ := jwtToken.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
	return tokenString
}

type TokenSet struct {
	ID                 primitive.ObjectID `bson:"_id" json:"-"`
	UserID             primitive.ObjectID `bson:"user_id" json:"-"`
	RefreshTokenID     string             `bson:"rt_id" json:"-"`
	PrevRefreshTokenID *string            `bson:"prt_id" json:"-"`
	UpdatedAt          time.Time          `bson:"updated_at" json:"-"`
}

func (tokenSet *TokenSet) MarshalBSON() ([]byte, error) {
	type TokenSetAlias TokenSet
	return bson.Marshal((*TokenSetAlias)(tokenSet))
}
