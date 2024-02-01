//go:build wireinject
// +build wireinject

//The build tag makes sure the stub is not built in the final build.

package tool

import (
	"github.com/google/wire"

	"github.com/bnb-chain/token-recover-approver/internal/config"
	"github.com/bnb-chain/token-recover-approver/internal/injection"
)

func Initialize(configPath string) (*Tool, error) {
	wire.Build(
		newTool,
		config.NewConfig,
		injection.InitLogger,
		injection.InitStore,
	)
	return &Tool{}, nil
}
