package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	WalletAddress string             `bson:"wallet_addresspackage models"`
	Email         string             `bson:"emailpackage models"`
	XLink         string             `bson:"x_link,omitempty"`
	Bio           string             `bson:"bio,omitempty"`
	CreatedAt     int64              `bson:"created_at"`
	UpdatedAt     int64              `bson:"updated_at"`
}
