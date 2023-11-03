package gorm

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	// database driver for gorm
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"

	"github.com/cosmos/cosmos-sdk/types"

	"github.com/bnb-chain/airdrop-service/internal/config"
	"github.com/bnb-chain/airdrop-service/internal/store"
)

type DataSourceTypeName string

const (
	Mysql      DataSourceTypeName = "mysql"
	Postgresql DataSourceTypeName = "postgres"
	Sqlite     DataSourceTypeName = "sqlite"
)

var _supportedDataSource = map[DataSourceTypeName]func(port uint, host, dbname, user, password, connectTimeout, readTimeout, writeTimeout string, sslmode bool) gorm.Dialector{
	Mysql: func(port uint, host, dbname, user, password, connectTimeout, readTimeout, writeTimeout string, sslmode bool) gorm.Dialector {
		return mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=UTC&time_zone=UTC&timeout=%s&readTimeout=%s&writeTimeout=%s", user, password, host, port, dbname, connectTimeout, readTimeout, writeTimeout))
	},
	Postgresql: func(port uint, host, dbname, user, password, connectTimeout, readTimeout, writeTimeout string, sslmode bool) gorm.Dialector {
		ssl := "disable"
		if sslmode {
			ssl = "allow"
		}
		return postgres.Open(fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s timezone=UTC", host, port, user, dbname, password, ssl))
	},
	Sqlite: func(port uint, host, dbname, user, password, connectTimeout, readTimeout, writeTimeout string, sslmode bool) gorm.Dialector {
		return sqlite.Open(fmt.Sprintf("%s.db", dbname))
	},
}

var _ store.Store = (*SQLStore)(nil)

func NewSQLStore(config *config.Config, options ...Option) (*SQLStore, error) {
	supported, ok := _supportedDataSource[DataSourceTypeName(config.Store.SqlStore.SQLDriver)]
	if !ok {
		return nil, fmt.Errorf("unsupported database driver: %s", config.Store.SqlStore.SQLDriver)
	}

	sqlDriver := supported(
		config.Store.SqlStore.Port, config.Store.SqlStore.Host, config.Store.SqlStore.DBName,
		config.Store.SqlStore.User, config.Store.SqlStore.Password,
		config.Store.SqlStore.ConnectTimeout, config.Store.SqlStore.ReadTimeout, config.Store.SqlStore.WriteTimeout,
		config.Store.SqlStore.SSLMode)

	engine, err := gorm.Open(sqlDriver, &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Store.SqlStore.DialTimeout)
	defer cancel()

	sql, err := engine.DB()
	if err != nil {
		return nil, err
	}
	err = sql.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	for _, opt := range options {
		err = opt.Apply(engine)
		if err != nil {
			return nil, err
		}
	}

	for _, migration := range Migrations {
		err := migration.Migrate(engine)
		if err != nil {
			return nil, err
		}
	}

	return &SQLStore{
		db: engine,
	}, nil
}

// SQLStore implements store.Store.
type SQLStore struct {
	db *gorm.DB
}

// GetAccountByAddress implements store.Store.
func (s *SQLStore) GetAccountByAddress(address types.AccAddress) (*store.Account, error) {
	var acc *Account
	result := s.db.Where("address = ?", address.String()).First(acc)
	if result.Error != nil {
		return nil, result.Error
	}

	return &store.Account{
		Address:       acc.Address,
		AccountNumber: acc.AccountNumber,
		SummaryCoins:  acc.SummaryCoins,
		Coins:         acc.Coins,
		FrozenCoins:   acc.FrozenCoins,
		LockedCoins:   acc.LockedCoins,
	}, nil
}

// GetAccountProofs implements store.Store.
func (s *SQLStore) GetAccountAssetProofs(address types.AccAddress, symbol string, tokenIndex int64) (proofs []string, err error) {
	var proof *Proof
	result := s.db.Where("address = ? AND index = ? AND denom = ?", address.String(), tokenIndex, symbol).First(proof)
	if result.Error != nil {
		return nil, result.Error
	}
	return proof.Proof, nil
}

// GetAssetBySymbol implements store.Store.
func (s *SQLStore) GetAssetBySymbol(symbol string) (asset *store.Asset, err error) {
	var assetModel *Asset
	result := s.db.Where("denom = ?", symbol).First(assetModel)
	if result.Error != nil {
		return nil, result.Error
	}
	return &store.Asset{
		Owner:  assetModel.Owner,
		Amount: assetModel.Amount,
	}, nil
}

// GetStateRoot implements store.Store.
func (s *SQLStore) GetStateRoot() (stateRoot string, err error) {
	var state *StateRoot
	result := s.db.First(state)
	if result.Error != nil {
		return "", result.Error
	}
	return state.StateRoot, nil
}
