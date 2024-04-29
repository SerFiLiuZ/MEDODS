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
	config *Config
}

func newServer(store store.Store, logger *utils.Logger, config *Config) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logger,
		store:  store,
		config: config,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/auth", s.authorizationUser()).Methods("GET")
	s.router.HandleFunc("/refresh", s.refreshUser()).Methods("POST")
}

func (s *server) authorizationUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guid := r.URL.Query().Get("guid")
		if guid == "" {
			s.logger.Errorf("Parameter 'guid' is required")
			http.Error(w, "Parameter 'guid' is required", http.StatusBadRequest)
			return
		}

		s.logger.Debugf("authorizationUser: guid: %v", guid)

		authUser, err := s.store.User().Receive(guid, s.config.JwtKey)
		if err != nil {
			s.logger.Errorf("authorizationUser: %v", err)
			http.Error(w, "authorizationUser: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "at",
			Value:    authUser.AtSigned,
			Expires:  authUser.AtExpiration,
			HttpOnly: true,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "rt",
			Value:    authUser.RtSigned,
			Expires:  authUser.RtExpiration,
			HttpOnly: true,
		})
	}
}

func (s *server) refreshUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
