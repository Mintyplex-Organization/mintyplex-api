package models

type DoTier1 struct {
	WalletAddress string `json:"wallet_address" bson:"_id"`
	Bio           string `json:"bio"`
	XLink         string `json:"x_link"`
}
