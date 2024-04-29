package sqlstore

import (
	"github.com/SerFiLiuZ/MEDODS/internal/app/models"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Receive(guid, jwtKey string) (*models.AuthUser, error) {
	au, err := r.store.tokenRepository.GetAccessRefreshTokens(guid, jwtKey)
	if err != nil {
		return nil, err
	}

	return au, nil
}
