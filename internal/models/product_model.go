package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddProduct struct {
	CoverImage  string    `bson:"image"`
	Name        string    `bson:"name"`
	Price       float64   `bson:"price"`
	Discount    float64   `bson:"discount"`
	Description string    `bson:"description"`
	Categories  string    `bson:"categories"`
	Quantity    int       `bson:"quantity"`
	Tags        []string  `bson:"tags"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserId      string             `bson:"user_id,omitempty"`
	CoverImage  string             `bson:"image"`
	Name        string             `bson:"name"`
	Price       float64            `bson:"price"`
	Discount    float64            `bson:"discount"`
	Description string             `bson:"description"`
	Categories  string             `bson:"categories"`
	Quantity    int                `bson:"quantity"`
	Tags        []string           `bson:"tags"`
	CreatedAt   int64              `bson:"created_at"`
	UpdatedAt   int64              `bson:"updated_at"`
}
