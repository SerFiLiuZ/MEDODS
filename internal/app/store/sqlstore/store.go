package sqlstore

import (
	"context"

	"github.com/SerFiLiuZ/MEDODS/internal/app/store"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	db              *mongo.Client
	ctx             *context.Context
	userRepository  *UserRepository
	tokenRepository *TokenRepository
}

func New(db *mongo.Client, ctx *context.Context) *Store {
	return &Store{
		db:  db,
		ctx: ctx,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

func (s *Store) Token() store.TokenRepository {
	if s.tokenRepository != nil {
		return s.tokenRepository
	}

	s.tokenRepository = &TokenRepository{
		store: s,
	}

	return s.tokenRepository
}
