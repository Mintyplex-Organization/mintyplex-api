package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	WalletAddress string             `bson:"_id,omitempty"`
	ID            primitive.ObjectID `bson:"uid,omitempty"`
	XLink         string             `bson:"x_link,omitempty"`
	Bio           string             `bson:"bio,omitempty"`
	Avatar        string             `bson:"avatar"`
	Products      string             `bson:"products"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}
