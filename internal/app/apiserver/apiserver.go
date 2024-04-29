package apiserver

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/SerFiLiuZ/MEDODS/internal/app/store/sqlstore"
	"github.com/SerFiLiuZ/MEDODS/internal/app/utils"
	"github.com/gorilla/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
	ctx    *context.Context
	err    error
}

func Start(config *Config, logger *utils.Logger) error {
	db := Connect(config.DBconnecturi)
	if db.err != nil {
		return db.err
	}
	defer func() {
		err := db.client.Disconnect(*db.ctx)
		if err != nil {
			return
		}
	}()

	logger.Infof("Connected to database")

	logger.Debugf("db: %v", db)

	store := sqlstore.New(db.client)

	logger.Debugf("store: %v", store)

	srv := newServer(store, logger, config)

	logger.Infof("Server started on port %s", config.Port)

	logger.Debugf("srv.store: %v", srv.store)

	return http.ListenAndServe(config.Port,
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With", "Cookie"}),
			handlers.ExposedHeaders([]string{"Set-Cookie"}),
			handlers.AllowCredentials(),
		)(srv))
}

func Connect(dbURI string) *MongoDB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		dbURI,
	))

	if err != nil {
		return &MongoDB{
			client: nil,
			ctx:    nil,
			err:    errors.New("Failed to connect DB: " + err.Error()),
		}
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return &MongoDB{
			client: nil,
			ctx:    nil,
			err:    errors.New("Failed to ping DB: " + err.Error()),
		}
	}

	return &MongoDB{
		client: client,
		ctx:    &ctx,
		err:    nil,
	}
}
