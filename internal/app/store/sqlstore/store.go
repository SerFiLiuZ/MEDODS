package sqlstore

import (
	"github.com/SerFiLiuZ/MEDODS/internal/app/store"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	db             *mongo.Client
	userRepository *UserRepository
}

func New(db *mongo.Client) *Store {
	return &Store{
		db: db,
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
