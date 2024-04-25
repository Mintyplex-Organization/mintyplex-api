package models

import "time"

type ReserveUsername struct {
	Email    string `json:"email" bson:"email"`
	Username string `json:"username" bson:"username"`
}

type Reserved struct {
	Email     string    `bson:"email"`
	Username  string    `bson:"username"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
