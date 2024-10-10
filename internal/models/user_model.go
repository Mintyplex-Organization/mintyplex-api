package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignUpUser struct {
	Name            string    `json:"name" bson:"name" binding:"required"`
	Email           string    `json:"email" bson:"email" binding:"required"`
	Password        string    `json:"password" bson:"password" binding:"required,min=8"`
	PasswordConfirm string    `json:"passwordConfirm" bson:"passwordConfirm,omitempty" binding:"required"`
	Provider        string    `json:"provider,omitempty" bson:"provider,omitempty"`
	Verified        bool      `json:"verified" bson:"verified"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}

type SignInInput struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

type UpdateDBUser struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string             `json:"name,omitempty" bson:"name,omitempty"`
	Email           string             `json:"email,omitempty" bson:"email,omitempty"`
	Password        string             `json:"password,omitempty" bson:"password,omitempty"`
	PasswordConfirm string             `json:"passwordConfirm,omitempty" bson:"passwordConfirm,omitempty"`
	Provider        string             `json:"provider" bson:"provider"`
	Verified        bool               `json:"verified,omitempty" bson:"verified,omitempty"`
	CreatedAt       time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt       time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Provider  string             `json:"provider" bson:"provider"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type DBResponse struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	Name            string             `json:"name" bson:"name"`
	Email           string             `json:"email" bson:"email"`
	Password        string             `json:"password" bson:"password"`
	PasswordConfirm string             `json:"passwordConfirm,omitempty" bson:"passwordConfirm,omitempty"`
	Verified        bool               `json:"verified" bson:"verified"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

func FilteredResponse(user *DBResponse) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

type DoTier1 struct {
	WalletAddress string `json:"wallet_address" bson:"_id" validate:"required"`
	XLink         string `json:"x_link" validate:"required"`
	Bio           string `json:"bio" validate:"required"`
	Avatar        string `json:"avatar" bson:"avatar"`
}

type User struct {
	WalletAddress string    `bson:"_id"`
	XLink         string    `bson:"x_link,omitempty"`
	Bio           string    `bson:"bio,omitempty"`
	Avatar        string    `bson:"avatar"`
	Products      []Product `json:"products"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at"`
}
