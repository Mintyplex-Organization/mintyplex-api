package models

type DoTier1 struct {
	Email         string `json:"email" validate:"required,email"`
	WalletAddress string `json:"wallet_address" validate:"required"`
	Bio           string `json:"bio" validate:"required"`
	XLink         string `json:"x_link" validate:"required"`
}
