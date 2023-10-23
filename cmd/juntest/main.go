package main

import (
	"github.com/Raitfolt/juntest/config"
	"github.com/Raitfolt/juntest/internal/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.LogPath)
	defer log.Sync()

	log.Info("Logger loaded")
}
