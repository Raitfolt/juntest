package main

import (
	"github.com/Raitfolt/juntest/config"
	"github.com/Raitfolt/juntest/internal/logger"
	"github.com/Raitfolt/juntest/internal/storage/psql"
	"go.uber.org/zap"
)

func main() {
	//CONFIG_PATH
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

}
