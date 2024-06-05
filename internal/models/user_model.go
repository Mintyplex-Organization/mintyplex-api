package models

import (
	"time"
)

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
