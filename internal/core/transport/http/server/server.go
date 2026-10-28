package server

import (
	"github.com/jackc/pgx/v5/pgxpool"
	swagger "github.com/swaggo/http-swagger/v2"
	"kontakt/internal/core/config"
	"kontakt/internal/core/transport/http/middleware"
	"kontakt/internal/core/transport/http/response"
	ct "kontakt/internal/features/catalog/transport"
	ot "kontakt/internal/features/orders/transport"
	"net/http"
	"time"
)

func New(c config.Config, db *pgxpool.Pool, cat *ct.Handler, orders *ot.Handler) *http.Server {
	m := http.NewServeMux()
	m.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		if e := db.Ping(r.Context()); e != nil {
			response.Error(w, 503, "database unavailable")
			return
		}
		response.JSON(w, 200, map[string]string{"status": "ok"})
	})
	m.HandleFunc("GET /api/v1/products", cat.List)
	m.HandleFunc("GET /api/v1/products/{id}", cat.Get)
	m.HandleFunc("POST /api/v1/products", middleware.Admin(c.AdminToken, cat.Save))
	m.HandleFunc("PUT /api/v1/products/{id}", middleware.Admin(c.AdminToken, cat.Save))
	m.HandleFunc("GET /api/v1/categories", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, 200, []string{"Смартфоны", "Ноутбуки", "Аудио", "Часы"})
	})
	m.HandleFunc("POST /api/v1/orders", orders.Create)
	m.HandleFunc("GET /api/v1/orders", middleware.Admin(c.AdminToken, orders.List))
	m.HandleFunc("PATCH /api/v1/orders/{id}/status", middleware.Admin(c.AdminToken, orders.Status))
	m.Handle("GET /swagger/", swagger.WrapHandler)
	return &http.Server{Addr: c.HTTPAddr, Handler: middleware.Wrap(c.CORSOrigin, m), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second}
}
