package injection

import (
	"github.com/rs/zerolog"

	"github.com/bnb-chain/airdrop-service/internal/config"
	"github.com/bnb-chain/airdrop-service/internal/version"
	"github.com/bnb-chain/airdrop-service/pkg/logger"
)

func InitLogger(config *config.Config) (*zerolog.Logger, error) {
	return logger.NewLogger(config.Logger.Level, config.Logger.Format, logger.WithStr("app_id", version.APPName))
}
