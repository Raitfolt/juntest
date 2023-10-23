package psql

import (
	"fmt"

	"github.com/Raitfolt/juntest/config"
	"go.uber.org/zap"
)

type Storage struct {
	//db *sql.DB
}

func New(log *zap.Logger, cfg *config.Config) {
	log.Info("create new storage connect")
	psqlInfo := fmt.Sprintf("host=%s port =%d user=%s password=%s sslmode=disalbe",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword)
	log.Info("postgresql", zap.String("connection string", psqlInfo))
}
