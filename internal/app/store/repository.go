package store

import (
	"github.com/SerFiLiuZ/MEDODS/internal/app/models"
)

type UserRepository interface {
	Receive(guid string, jwtKey []byte) (*models.AuthUser, error)
}

type TokenRepository interface {
	GetAccessRefreshTokens(guid string, jwtKey []byte) (*models.AuthUser, error)
	GetHashedToken(token string) (string, error)
}
