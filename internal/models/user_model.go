package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type User struct {
// 	ID            primitive.ObjectID `bson:"_id,omitempty"`
// 	WalletAddress string             `bson:"wallet_address,required"`
// 	Email         string             `bson:"email,required"`
// 	XLink         string             `bson:"x_link,omitempty"`
// 	Bio           string             `bson:"bio,omitempty"`
// 	CreatedAt     int64              `bson:"created_at"`
// 	UpdatedAt     int64              `bson:"updated_at"`
// }

type User struct {
	WalletAddress string             `bson:"_id,omitempty"` // Now the primary identifier
	ID            primitive.ObjectID `bson:"id,omitempty"`  // Moved down
	Email         string             `bson:"email,required"`
	XLink         string             `bson:"x_link,omitempty"`
	Bio           string             `bson:"bio,omitempty"`
	Avatar        string             `bson:"avatar"`
	CreatedAt     int64              `bson:"created_at"`
	UpdatedAt     int64              `bson:"updated_at"`
}
