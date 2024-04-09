package models

// type DoTier1 struct {
// 	Email         string `json:"email"`
// 	WalletAddress string `json:"wallet_address" `
// 	Bio           string `json:"bio" `
// 	XLink         string `json:"x_link" `
// }

type DoTier1 struct {
	WalletAddress string `json:"wallet_address" bson:"_id,omitempty"` // Now primary identifier
	Email         string `json:"email"`
	Bio           string `json:"bio"`
	XLink         string `json:"x_link"`
}
