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

	s.router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
	})

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
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		userInfo, err := s.authService.ExchangeAuthCode(r.Context(), req.ServerAuthCode, req.RedirectURI)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		u := &models.User{
			GoogleID:     userInfo.GoogleID,
			RefreshToken: userInfo.RefreshToken,
			AccessToken:  userInfo.AccessToken,
		}

		if err := s.store.User().FindOrCreateUser(u); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		JWT, err := jwt.GenerateJWT(u.ID)
		if err != nil {
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
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		u.Sanitize()
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

		entries, err := s.store.Sub().GetAllSubsForUser(userID, r.URL.Query().Get("category"))
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

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
			s.error(w, r, http.StatusUnauthorized, errors.New("no token"))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := jwt.ParseJWT(tokenString)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		allowedOrigins := map[string]bool{
			"http://localhost:3000": true,
			"http://localhost:8081": true,
			"http://127.0.0.1:3000": true,
			"http://192.168.0.0/16": true,
			"http://10.0.0.0/8":     true,
		}

		if strings.HasPrefix(origin, "http://192.168.") ||
			strings.HasPrefix(origin, "http://10.") ||
			strings.HasPrefix(origin, "http://172.") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		w.Header().Set("Access-Control-Expose-Headers", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
