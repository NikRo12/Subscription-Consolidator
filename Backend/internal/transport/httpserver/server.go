package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"
	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/services/email"
	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/services/jwt"
	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/services/task"
	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/store"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type contextKey string

const userIDKey contextKey = "user_id"

type server struct {
	router      *mux.Router
	logger      *logrus.Logger
	store       store.Store
	authService *email.AuthService
	taskService *task.TaskService
}

func newServer(store store.Store, logger *logrus.Logger, redisClient *redis.Client, authService *email.AuthService) *server {
	s := &server{
		router:      mux.NewRouter(),
		logger:      logger,
		store:       store,
		authService: authService,
		taskService: task.NewTaskService(*redisClient),
	}

	s.configureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(corsMiddleware)
	s.router.HandleFunc("/auth/google", s.handleGoogleAuth()).Methods("POST")
	s.router.HandleFunc("/subscriptions", s.authenticateUser(s.handleSubscriptions())).Methods("GET")
}

func (s *server) handleGoogleAuth() http.HandlerFunc {
	type request struct {
		ServerAuthCode string `json:"serverAuthCode"`
		RedirectURI    string `json:"redirectUri"`
	}

	type response struct {
		JWT string `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("POST /auth/google: request received")

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.logger.Errorf("POST /auth/google: failed to decode request body: %v", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.logger.Infof("POST /auth/google: exchanging auth code, redirectURI=%q", req.RedirectURI)

		userInfo, err := s.authService.ExchangeAuthCode(r.Context(), req.ServerAuthCode, req.RedirectURI)
		if err != nil {
			s.logger.Errorf("POST /auth/google: failed to exchange auth code: %v", err)
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		s.logger.Infof("POST /auth/google: auth code exchanged, googleID=%q", userInfo.GoogleID)

		u := &models.User{
			GoogleID:     userInfo.GoogleID,
			RefreshToken: userInfo.RefreshToken,
			AccessToken:  userInfo.AccessToken,
		}

		if err := s.store.User().FindOrCreateUser(u); err != nil {
			s.logger.Errorf("POST /auth/google: failed to find or create user googleID=%q: %v", u.GoogleID, err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.logger.Infof("POST /auth/google: user found/created id=%d googleID=%q", u.ID, u.GoogleID)

		JWT, err := jwt.GenerateJWT(u.ID)
		if err != nil {
			s.logger.Errorf("POST /auth/google: failed to generate JWT for userID=%d: %v", u.ID, err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		task := models.Task{
			UserID:       u.ID,
			RefreshToken: u.RefreshToken,
			AccessToken:  u.AccessToken,
			MessageID:    0,
		}

		if err := s.taskService.SendTask(r.Context(), &task); err != nil {
			s.logger.Errorf("POST /auth/google: failed to send task to Redis for userID=%d: %v", u.ID, err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.logger.Infof("POST /auth/google: task sent to Redis for userID=%d", u.ID)

		u.Sanitize()
		s.logger.Infof("POST /auth/google: responding 200 for userID=%d", u.ID)
		s.respond(w, r, http.StatusOK, response{JWT: JWT})
	}
}

func (s *server) handleSubscriptions() http.HandlerFunc {
	type currencyTotal struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
	}

	type response struct {
		MonthlySpend []*currencyTotal `json:"monthly_spend"`
		Items        []*models.Entry  `json:"items"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(userIDKey).(int)
		category := r.URL.Query().Get("category")

		s.logger.Infof("GET /subscriptions: userID=%d category=%q", userID, category)

		entries, err := s.store.Sub().GetAllSubsForUser(userID, category)
		if err != nil {
			s.logger.Errorf("GET /subscriptions: failed to get subs for userID=%d: %v", userID, err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.logger.Infof("GET /subscriptions: found %d entries for userID=%d", len(entries), userID)

		totalsMap := make(map[string]float64)
		for _, e := range entries {
			if !e.IsActive {
				continue
			}
			price, _ := e.Price.Float64()
			if e.Period == models.Yearly {
				price = price / 12
			}
			totalsMap[e.Currency] += price
		}

		totals := make([]*currencyTotal, 0, len(totalsMap))
		for currency, amount := range totalsMap {
			totals = append(totals, &currencyTotal{
				Amount:   amount,
				Currency: currency,
			})
		}

		s.respond(w, r, http.StatusOK, response{
			MonthlySpend: totals,
			Items:        entries,
		})
	}
}

func (s *server) authenticateUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.logger.Warn("authenticateUser: missing Authorization header")
			s.error(w, r, http.StatusUnauthorized, errors.New("no token"))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := jwt.ParseJWT(tokenString)
		if err != nil {
			s.logger.Errorf("authenticateUser: invalid token: %v", err)
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		s.logger.Infof("authenticateUser: authenticated userID=%d", userID)

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.logger.Errorf("HTTP %d %s %s: %v", code, r.Method, r.URL.Path, err)
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
