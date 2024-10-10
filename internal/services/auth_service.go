package services

import "mintyplex-api/internal/models"

type AuthService interface {
	SignUpUser(*models.SignUpUser) (*models.UserResponse, error)
	SignInUser(*models.SignInInput) (*models.UserResponse, error)
}