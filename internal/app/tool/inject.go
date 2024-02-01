package tool

import (
	"github.com/rs/zerolog"

	"github.com/bnb-chain/token-recover-approver/internal/config"
	"github.com/bnb-chain/token-recover-approver/internal/store"
)

type Tool struct {
	logger *zerolog.Logger
	config *config.Config

	store store.Store
}

func newTool(
	logger *zerolog.Logger,
	config *config.Config,
	store store.Store,
) *Tool {
	return &Tool{
		logger: logger,
		config: config,
		store:  store,
	}
}
