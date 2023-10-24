package main

import (
	"net/http"

	"github.com/Raitfolt/juntest/config"
	"github.com/Raitfolt/juntest/internal/httpServer/handlers/add"
	"github.com/Raitfolt/juntest/internal/logger"
	"github.com/Raitfolt/juntest/internal/storage/psql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.LogPath)
	log.Info("logger loaded")
	defer log.Sync()

	db, err := psql.New(log, cfg)
	if err != nil {
		log.Fatal("postgre connection", zap.String("error", err.Error()))
	}
	defer db.DB.Close()

	log.Info("database loaded")

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)

	router.Post("/add", add.New(log, db))

	address := cfg.Host + ":" + cfg.Port
	srv := &http.Server{
		Addr:    address,
		Handler: router,
	}

	log.Info("starting server", zap.String("address", address))
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("server failed", zap.Error(err))
	}
}
