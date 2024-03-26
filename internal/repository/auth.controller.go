package repository

import (
	"context"
	"errors"
	"mintyplex-api/internal/models"
	"mintyplex-api/internal/utils"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthRepositoryImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewAuthRepository(collection *mongo.Collection, ctx context.Context) AuthRepository {
	return &AuthRepositoryImpl{collection, ctx}
}

func (asi *AuthRepositoryImpl) SignUpUser(user *models.SignUpReq) (*models.ReqResponse, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Email = strings.ToLower(user.Email)
	user.PasswordConfirm = ""

	encryptedPassword, _ := utils.EncryptPassword(user.Password)
	user.Password = encryptedPassword
	res, err := asi.collection.InsertOne(asi.ctx, &user)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("email already exists")
		}
		return nil, err
	}

	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: opt}

	if _, err := asi.collection.Indexes().CreateOne(asi.ctx, index); err != nil {
		return nil, errors.New("error creating index for email")
	}

	var newUser *models.ReqResponse
	query := bson.M{"_id": res.InsertedID}

	err = asi.collection.FindOne(asi.ctx, query).Decode(&newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (asi *AuthRepositoryImpl) LoginUser(user *models.LoginReq) (*models.ReqResponse, error) {
	return nil, nil
}
