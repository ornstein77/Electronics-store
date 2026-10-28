package main

import (
	"context"
	"errors"
	"go.uber.org/zap"
	_ "kontakt/docs"
	"kontakt/internal/core/config"
	"kontakt/internal/core/logger"
	"kontakt/internal/core/repository/postgres/pool"
	"kontakt/internal/core/transport/http/server"
	cr "kontakt/internal/features/catalog/repository"
	cs "kontakt/internal/features/catalog/service"
	ct "kontakt/internal/features/catalog/transport"
	or "kontakt/internal/features/orders/repository"
	oservice "kontakt/internal/features/orders/service"
	ot "kontakt/internal/features/orders/transport"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Kontakt API
// @version 1.0
// @description Electronics catalog and orders API
// @BasePath /api/v1
// @securityDefinitions.apikey AdminToken
// @in header
// @name Authorization
func main() {
	if e := run(); e != nil {
		zap.L().Error("application stopped", zap.Error(e))
		os.Exit(1)
	}
}
func run() error {
	l, e := logger.New()
	if e != nil {
		return e
	}
	defer l.Sync()
	zap.ReplaceGlobals(l)
	c, e := config.Load()
	if e != nil {
		return e
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	connect, cancel := context.WithTimeout(ctx, 10*time.Second)
	db, e := pool.New(connect, c.DatabaseURL)
	cancel()
	if e != nil {
		return e
	}
	defer db.Close()
	s := server.New(c, db, &ct.Handler{Service: &cs.Service{Repo: &cr.Repository{DB: db}}}, &ot.Handler{Service: &oservice.Service{Repo: &or.Repository{DB: db}}})
	ch := make(chan error, 1)
	go func() { l.Info("listening", zap.String("address", c.HTTPAddr)); ch <- s.ListenAndServe() }()
	select {
	case e = <-ch:
		if errors.Is(e, http.ErrServerClosed) {
			return nil
		}
		return e
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return s.Shutdown(shutdown)
	}
}
