package models

import "time"

type UserProfileResponse struct {
	WalletAddress string    `json:"wallet_address" bson:"_id,omitempty"`
	Bio           string    `json:"bio"`
	XLink         string    `json:"x_link"`
	Avatar        string    `json:"avatar"`
	Products      []Product `json:"products"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
