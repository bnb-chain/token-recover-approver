//go:build wireinject
// +build wireinject

//The build tag makes sure the stub is not built in the final build.

package app

import (
	"github.com/google/wire"

	"github.com/bnb-chain/airdrop-service/internal/config"
	"github.com/bnb-chain/airdrop-service/internal/module/approval"
	"github.com/bnb-chain/airdrop-service/internal/module/http"
	"github.com/bnb-chain/airdrop-service/internal/wireset"
)

func Initialize(configPath string) (Application, error) {
	wire.Build(
		newApplication,
		config.NewConfig,
		wireset.InitLogger,
		wireset.InitKeyManager,
		wireset.InitStore,
		approval.NewApprovalService,
		http.NewHttpServer,
	)
	return Application{}, nil
}
