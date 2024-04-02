package models

type UserProfile struct {
	ID            string `json:"id"`
	WalletAddress string `json:"wallet_address"`
	Email         string `json:"email"`
	Avatar        string `json:"avatar"`
	Bio           string `json:"bio"`
	XLink         string `json:"x_link"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
}