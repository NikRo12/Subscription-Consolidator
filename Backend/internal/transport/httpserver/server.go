package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

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
	s.router.HandleFunc("/subscriptions", s.authenticateUser(s.handleSubscriptions())).Methods("Get")
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

		u := &models.User{
			RefreshToken: req.RefreshToken,
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

func (s *server) handleSubscriptions() http.HandlerFunc {
	type response struct{}

	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("user_id").(int)

		entries, err := s.store.Sub().GetAllSubsForUser(userID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, entries)
	}
}

func (s *server) authenticateUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.error(w, r, http.StatusUnauthorized, errors.New("no token"))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := services.ParseJWT(tokenString)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next(w, r.WithContext(ctx))
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
