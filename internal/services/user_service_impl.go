package services

import (
	"context"
	"errors"
	"mintyplex-api/internal/models"
	"mintyplex-api/internal/utils"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserServiceImpl(collection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{collection, ctx}
}

func (us *UserServiceImpl) GetUserById(id string) (*models.UserResponse, error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var user *models.UserResponse

	query := bson.M{"_id": oid}
	err := us.collection.FindOne(us.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.UserResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) GetUserByEmail(email string) (*models.UserResponse, error) {
	var user *models.UserResponse

	query := bson.M{"email": strings.ToLower(email)}
	err := us.collection.FindOne(us.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.UserResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) UpsertUser(email string, data *models.UpdateDBUser) (*models.UserResponse, error) {
	doc, err := utils.ToBsonDoc(data)
	if err != nil {
		return nil, err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	query := bson.D{{Key: "email", Value: email}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := us.collection.FindOneAndUpdate(us.ctx, query, update, opts)

	var updatedPost *models.UserResponse

	if err := res.Decode(&updatedPost); err != nil {
		return nil, errors.New("no post with that Id exists")
	}

	return updatedPost, nil
}
