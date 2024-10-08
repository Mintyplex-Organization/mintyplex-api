package services

import "mintyplex-api/internal/models"

type UserService interface{
	GetUserById(string) (*models.UserResponse, error)
	GetUserByEmail(string) (*models.UserResponse, error)
	UpsertUser(string, *models.UpdateDBUser) (*models.UserResponse, error)
}