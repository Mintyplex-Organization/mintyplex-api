package repository

import "mintyplex-api/internal/models"

type UserRepository interface {
	FindUserById(string) (*models.ReqResponse, error)
	FindUserByEmail(string) (*models.ReqResponse, error)
}