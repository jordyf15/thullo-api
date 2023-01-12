package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MemberRole string

const (
	MemberRoleMember = "member"
	MemberRoleAdmin  = "admin"
)

type BoardMember struct {
	ID      primitive.ObjectID `json:"id"`
	UserID  primitive.ObjectID `json:"user_id"`
	BoardID primitive.ObjectID `json:"board_id"`
	Role    MemberRole         `json:"role"`
}
