package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignUpReq struct {
	Name            string `json:"name" bson:"name" binding:"required"`
	Email           string
	Password        string    `json:"password" bson:"password" binding:"required,min=8"`
	PasswordConfirm string    `json:"passwordConfirm" bson:"passwordConfirm,omitempty" binding:"required"`
	ImageUrl        string    `json:"imageUrl" bson:"imageUrl"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}

type LoginReq struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password"`
}

type ReqResponse struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	Name            string             `json:"name" bson:"name"`
	Email           string             `json:"email" bson:"email"`
	Password        string             `json:"password" bson:"password"`
	PasswordConfirm string             `json:"passwordConfirm,omitempty" bson:"passwordConfirm,omitempty"`
	ImageUrl        string    `json:"imageUrl" bson:"imageUrl"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	ImageUrl        string    `json:"imageUrl" bson:"imageUrl"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func FilteredResponse(user *ReqResponse) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		ImageUrl: user.ImageUrl,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
