package store

import (
	"github.com/SerFiLiuZ/MEDODS/internal/app/models"
)

type UserRepository interface {
	Receive(guid, jwtKey string) (*models.AuthUser, error)
}

type TokenRepository interface {
	GetAccessRefreshTokens(guid, jwtKey string) (*models.AuthUser, error)
}
