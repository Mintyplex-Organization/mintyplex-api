package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Email        string             `bson:"email,required"`
	PasswordHash string             `bson:"password_hash"`
	Tokens       []string           `bson:"tokens"`
	CreatedAt    int64              `bson:"created_at"`
	UpdatedAt    int64              `bson:"updated_at"`
}

type UserProfileResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type SignUp struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Mobile    string `json:"mobile" validate:"numeric,len=10"`
	Password  string `json:"password" validate:"required,min=8,max=20"`
}

type SignIn struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}