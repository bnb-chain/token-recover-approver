//go:build wireinject
// +build wireinject

//The build tag makes sure the stub is not built in the final build.

package http

import (
	"github.com/google/wire"

	"github.com/bnb-chain/airdrop-service/internal/config"
	"github.com/bnb-chain/airdrop-service/internal/delivery/approval"
	"github.com/bnb-chain/airdrop-service/internal/delivery/http"
	"github.com/bnb-chain/airdrop-service/internal/wireset"
)

func Initialize(configPath string) (Application, error) {
	wire.Build(
		newApplication,
		config.NewConfig,
		approval.NewApprovalService,
		wireset.InitLogger,
		wireset.InitKeyManager,
		http.NewHttpServer,
	)
	return Application{}, nil
}
