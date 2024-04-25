package models

import "time"

// type UserProfile struct {
// 	ID            string `json:"id"`
// 	WalletAddress string `json:"wallet_address"`
// 	Avatar        string `json:"avatar"`
// 	Bio           string `json:"bio"`
// 	XLink         string `json:"x_link"`
// 	CreatedAt     int64  `json:"created_at"`
// 	UpdatedAt     int64  `json:"updated_at"`
// }

type UserProfile struct {
	WalletAddress string    `json:"wallet_address" bson:"_id,omitempty"`
	ID            string    `json:"uid"`
	Avatar        string    `json:"avatar"`
	Bio           string    `json:"bio"`
	XLink         string    `json:"x_link"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
