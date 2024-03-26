package repository

import "mintyplex-api/internal/models"

type AuthRepository interface {
	SignUpUser(*models.SignUpReq) (*models.ReqResponse, error)
	LoginUser(*models.LoginReq) (*models.ReqResponse, error)
}

