package wireset

import (
	"errors"

	"github.com/bnb-chain/airdrop-service/internal/common"
	"github.com/bnb-chain/airdrop-service/internal/config"
	"github.com/bnb-chain/airdrop-service/internal/store"
	"github.com/bnb-chain/airdrop-service/internal/store/goleveldb"
	"github.com/bnb-chain/airdrop-service/internal/store/gorm"
	"github.com/bnb-chain/airdrop-service/internal/store/memory"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog"
)

type StoreType string

const (
	MemoryStore  StoreType = "memory"
	LevelDBStore StoreType = "leveldb"
	GORMStore    StoreType = "gorm"
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
			config.ChainID,
			config.Store.MemoryStore.StateRoot,
			config.Store.MemoryStore.Assets,
			config.Store.MemoryStore.Accounts,
			config.Store.MemoryStore.MerkleProofs,
		)
	case LevelDBStore:
		// TODO: implement
		return goleveldb.NewKVStore()
	case GORMStore:
		// TODO: implement
		return gorm.NewSQLStore()
	default:
		return nil, errors.New("invalid store type")
	}
}
