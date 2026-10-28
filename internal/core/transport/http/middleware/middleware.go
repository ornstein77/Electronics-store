package middleware

import (
	"crypto/subtle"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"kontakt/internal/core/transport/http/response"
	"net/http"
	"strings"
	"time"
)

func Admin(token string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") || subtle.ConstantTimeCompare([]byte(v), []byte(token)) != 1 {
			response.Error(w, 401, "admin authorization required")
			return
		}
		next(w, r)
	}
}
func Wrap(origin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		id := uuid.NewString()
		w.Header().Set("X-Request-ID", id)
		defer func() {
			if v := recover(); v != nil {
				zap.L().Error("panic", zap.Any("error", v))
				response.Error(w, 500, "internal server error")
			}
			zap.L().Info("request", zap.String("request_id", id), zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Duration("duration", time.Since(start)))
		}()
		if r.Header.Get("Origin") == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, OPTIONS")
		}
		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}
		next.ServeHTTP(w, r)
	})
}
