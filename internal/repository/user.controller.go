package repository

import (
	"context"
	"mintyplex-api/internal/models"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserServiceImpl(collection *mongo.Collection, ctx context.Context) UserRepository {
	return &UserServiceImpl{collection, ctx}
}

func (us *UserServiceImpl) FindUserById(id string) (*models.ReqResponse, error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var user *models.ReqResponse

	query := bson.M{"_id": oid}
	err := us.collection.FindOne(us.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.ReqResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) FindUserByEmail(email string) (*models.ReqResponse, error) {
	var user *models.ReqResponse

	query := bson.M{"email": strings.ToLower(email)}
	err := us.collection.FindOne(us.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.ReqResponse{}, err
		}
		return nil, err
	}

	return user, nil
}
