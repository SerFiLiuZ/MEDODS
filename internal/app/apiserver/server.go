package apiserver

import (
	"net/http"

	"github.com/SerFiLiuZ/MEDODS/internal/app/store"
	"github.com/SerFiLiuZ/MEDODS/internal/app/utils"
	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
	logger *utils.Logger
	store  store.Store
}

func newServer(store store.Store, logger *utils.Logger) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logger,
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	//
}
