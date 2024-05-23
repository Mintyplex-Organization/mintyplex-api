package models

import (
	"time"
)

type User struct {
	WalletAddress string    `bson:"_id"`
	XLink         string    `bson:"x_link,omitempty"`
	Bio           string    `bson:"bio,omitempty"`
	Avatar        string    `bson:"avatar"`
	Products      []Product `json:"products"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at"`
}
