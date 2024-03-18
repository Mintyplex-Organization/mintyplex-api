package repository

import "mintyplex-api/internal/models"

type AuthRepository interface {
	SignUpReq(*models.SignUpReq) (*models.ReqResponse, error)
	LoginReq(*models.LoginReq) (*models.ReqResponse, error)
}

