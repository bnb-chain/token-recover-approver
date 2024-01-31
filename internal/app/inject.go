package app

import (
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/bnb-chain/token-recover-approver/internal/config"
	"github.com/bnb-chain/token-recover-approver/internal/module/http"
	"github.com/bnb-chain/token-recover-approver/internal/version"

	"github.com/rs/zerolog"
)

type Application struct {
	logger     *zerolog.Logger
	config     *config.Config
	httpServer *http.HttpServer
}

func (application Application) Start() error {
	application.logger.Info().Str("app_version", version.AppVersion).Str("git_commit", version.GitCommit).Str("git_commit_date", version.GitCommitDate).Msg("version info")
	eg := errgroup.Group{}
	eg.Go(func() error {
		application.logger.Info().Msgf("http server listen %s:%d", application.config.HTTP.Addr, application.config.HTTP.Port)
		return application.httpServer.Run(fmt.Sprintf("%s:%d", application.config.HTTP.Addr, application.config.HTTP.Port))
	})
	eg.Go(func() error {
		if !application.config.Metrics.Enable {
			return nil
		}
		application.logger.Info().Msgf("metrics server listen %s:%d", application.config.HTTP.Addr, application.config.HTTP.Port)
		return application.httpServer.RunMetrics(fmt.Sprintf("%s:%d", application.config.Metrics.Addr, application.config.Metrics.Port), application.config.Metrics.Path, application.config.Metrics.PProf)
	})

	return eg.Wait()
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
