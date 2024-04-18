package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type User struct {
// 	ID            primitive.ObjectID `bson:"_id,omitempty"`
// 	WalletAddress string             `bson:"wallet_address,required"`
// 	XLink         string             `bson:"x_link,omitempty"`
// 	Bio           string             `bson:"bio,omitempty"`
// 	CreatedAt     int64              `bson:"created_at"`
// 	UpdatedAt     int64              `bson:"updated_at"`
// }

type User struct {
	WalletAddress string             `bson:"_id,omitempty"` 
	ID            primitive.ObjectID `bson:"uid,omitempty"` 
	XLink         string             `bson:"x_link,omitempty"`
	Bio           string             `bson:"bio,omitempty"`
	Avatar        string             `bson:"avatar"`
	CreatedAt     int64              `bson:"created_at"`
	UpdatedAt     int64              `bson:"updated_at"`
}
