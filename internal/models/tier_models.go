package models

type DoTier1 struct {
	Avatar        string `json:"avatar" bson:"avatar"`
	WalletAddress string `json:"wallet_address" bson:"_id"`
	Bio           string `json:"bio"`
	XLink         string `json:"x_link"`
}
