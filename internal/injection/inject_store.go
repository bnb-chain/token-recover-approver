package injection

import (
	"errors"

	"github.com/rs/zerolog"
	gormLogger "gorm.io/gorm/logger"

	"github.com/cosmos/cosmos-sdk/types"

	"github.com/bnb-chain/token-recover-approver/internal/common"
	"github.com/bnb-chain/token-recover-approver/internal/config"
	"github.com/bnb-chain/token-recover-approver/internal/store"
	"github.com/bnb-chain/token-recover-approver/internal/store/gorm"
	"github.com/bnb-chain/token-recover-approver/internal/store/memory"
)

type StoreType string

const (
	MemoryStore StoreType = "memory"
	GORMStore   StoreType = "gorm"
)

func initSDK(config *config.Config, logger *zerolog.Logger) {
	logger.Info().Str("chain_id", config.ChainID).Msg("init sdk config")
	sdkConfig := types.GetConfig()
	sdkConfig.SetBech32PrefixForAccount("bnb", "bnbp")
	if config.ChainID != common.MainnetChainID {
		sdkConfig.SetBech32PrefixForAccount("tbnb", "bnbp")
		logger.Debug().Str("chain_id", config.ChainID).Msg("set bech32 prefix to tbnb")
	}
}

func InitStore(config *config.Config, logger *zerolog.Logger) (store.Store, error) {
	initSDK(config, logger)
	switch StoreType(config.Store.Driver) {
	case MemoryStore:
		return memory.NewMemoryStore(
			config.Store.MemoryStore.Accounts,
			config.Store.MemoryStore.MerkleProofs,
		)
	case GORMStore:
		return gorm.NewSQLStore(
			config,
			gorm.SetConnMaxLifetime(config.Store.SqlStore.MaxLifetime),
			gorm.SetConnMaxIdleTime(config.Store.SqlStore.MaxIdleTime),
			gorm.SetMaxIdleConns(config.Store.SqlStore.MaxIdleConn),
			gorm.SetMaxOpenConns(config.Store.SqlStore.MaxOpenConn),
			gorm.SetLogLevel(gormLogger.LogLevel(config.Store.SqlStore.LogLevel)),
		)
	default:
		return nil, errors.New("invalid store type")
	}
}
