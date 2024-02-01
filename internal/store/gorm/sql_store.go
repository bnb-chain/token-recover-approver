package gorm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	// database driver for gorm
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"

	"github.com/cosmos/cosmos-sdk/types"

	"github.com/bnb-chain/token-recover-approver/internal/config"
	"github.com/bnb-chain/token-recover-approver/internal/store"
	"github.com/bnb-chain/token-recover-approver/pkg/util"
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

// GetAccountProof implements store.Store.
func (s *SQLStore) GetAccountAssetProof(address types.AccAddress, symbol string) (*store.Proof, error) {
	var proof Proof
	result := s.db.Where("address = ? AND denom = ?", util.EncodeBytesToHex(address), symbol).First(&proof)
	if result.Error != nil {
		return nil, result.Error
	}

	return &store.Proof{
		Address: types.AccAddress(util.MustDecodeHexToBytes(proof.Address)),
		Denom:   proof.Denom,
		Amount:  proof.Amount,
		Proof:   util.MustDecodeHexArrayToBytes(strings.Split(proof.Proof, ",")),
	}, nil
}

// InsertAccountProof implements a function to insert account proof.
func (s *SQLStore) InsertAccountAssetProof(proof *store.Proof) error {
	dbProof := &Proof{
		Address: util.EncodeBytesToHex(proof.Address[:]),
		Denom:   proof.Denom,
		Amount:  proof.Amount,
		Proof:   strings.Join(util.EncodeBytesArrayToHex(proof.Proof), ","),
	}

	result := s.db.Create(dbProof)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Close implements store.Store.
func (s *SQLStore) Close() error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
