package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"
	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/services"
	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func newServer(store store.Store, logger *logrus.Logger) *server {
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
	s.router.HandleFunc("/auth/google", s.handleGoogleAuth()).Methods("Post")
}

func (s *server) handleGoogleAuth() http.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refresh_token"`
	}

	type response struct {
		JWT  string      `json:"token"`
		User models.User `json:"user"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		//googleUser, err := services.ExchangeCodeWithGoogle(req.RefreshToken)
		// if err != nil {
		// 	s.error(w, r, http.StatusUnauthorized, err)
		// 	return
		// }

		u := &models.User{
			Email:        "user@example.ru",
			RefreshToken: "123456789",
		}

		if err := s.store.User().FindOrCreateUser(u); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		JWT, err := services.GenerateJWT(u.ID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusOK, response{JWT: JWT, User: *u})
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
