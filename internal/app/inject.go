package app

import (
	"golang.org/x/sync/errgroup"

	"github.com/bnb-chain/token-recover-approver/internal/config"
	"github.com/bnb-chain/token-recover-approver/internal/module/http"
	"github.com/bnb-chain/token-recover-approver/internal/store"
	"github.com/bnb-chain/token-recover-approver/internal/version"

	"github.com/rs/zerolog"
)

type Application struct {
	logger     *zerolog.Logger
	config     *config.Config
	httpServer *http.HttpServer
	store      store.Store
}

func (application Application) Start() error {
	application.logger.Info().Str("app_version", version.AppVersion).Str("git_commit", version.GitCommit).Str("git_commit_date", version.GitCommitDate).Msg("version info")
	eg := errgroup.Group{}
	eg.Go(func() error {
		application.logger.Info().Msgf("http server listen %s:%d", application.config.HTTP.Addr, application.config.HTTP.Port)
		return application.httpServer.Run(application.config.HTTP)
	})
	eg.Go(func() error {
		if !application.config.Metrics.Enable {
			return nil
		}
		application.logger.Info().Msgf("metrics server listen %s:%d", application.config.Metrics.Addr, application.config.Metrics.Port)
		return application.httpServer.RunMetrics(application.config.Metrics)
	})

	return eg.Wait()
}

func (application Application) Stop() error {
	application.logger.Info().Msg("shutdown http server ...")
	if err := application.httpServer.Shutdown(); err != nil {
		return err
	}
	application.logger.Info().Msg("http server is closed")

	application.logger.Info().Msg("shutdown store ...")
	if err := application.store.Close(); err != nil {
		return err
	}
	application.logger.Info().Msg("store is closed")
	return nil
}

func newApplication(
	logger *zerolog.Logger,
	config *config.Config,
	httpServer *http.HttpServer,
	store store.Store,
) Application {
	return Application{
		logger:     logger,
		config:     config,
		httpServer: httpServer,
		store:      store,
	}
}
