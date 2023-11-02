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

	return &SQLStore{
		db: engine,
	}, nil
}

// SQLStore implements store.Store.
type SQLStore struct {
	db *gorm.DB
}

// GetAccountByAddress implements store.Store.
func (*SQLStore) GetAccountByAddress(address types.AccAddress) (account *store.Account, err error) {
	panic("unimplemented")
}

// GetAccountProofs implements store.Store.
func (*SQLStore) GetAccountAssetProofs(address types.AccAddress, symbol string, tokenIndex int64) (proofs []string, err error) {
	panic("unimplemented")
}

// GetAssetBySymbol implements store.Store.
func (*SQLStore) GetAssetBySymbol(symbol string) (asset *store.Asset, err error) {
	panic("unimplemented")
}

// GetStateRoot implements store.Store.
func (*SQLStore) GetStateRoot() (stateRoot string, err error) {
	panic("unimplemented")
}
