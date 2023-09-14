package http

import (
	"fmt"

	"github.com/bnb-chain/airdrop-service/internal/config"
	"github.com/bnb-chain/airdrop-service/internal/delivery/http"
	"github.com/bnb-chain/airdrop-service/internal/version"

	"github.com/rs/zerolog"
)

type Application struct {
	logger     *zerolog.Logger
	config     *config.Config
	httpServer *http.HttpServer
}

func (application Application) Start() error {
	application.logger.Info().Str("app_version", version.AppVersion).Str("git_commit", version.GitCommit).Str("git_commit_date", version.GitCommitDate).Msg("version info")
	application.logger.Info().Msgf("http server listen %s:%d", application.config.HTTP.Addr, application.config.HTTP.Port)
	return application.httpServer.Run(fmt.Sprintf("%s:%d", application.config.HTTP.Addr, application.config.HTTP.Port))
}

func (application Application) Stop() error {
	application.logger.Info().Msg("shutdown http server ...")
	defer application.logger.Info().Msg("http server is closed")
	return application.httpServer.Shutdown()
}

func newApplication(
	logger *zerolog.Logger,
	config *config.Config,
	httpServer *http.HttpServer,
) Application {
	return Application{
		logger:     logger,
		config:     config,
		httpServer: httpServer,
	}
}
