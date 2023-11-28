package config

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	gormLogger "gorm.io/gorm/logger"

	"github.com/bnb-chain/token-recover-approver/pkg/logger"
)

type Config struct {
	ChainID          string       `mapstructure:"chain_id"`
	MerkleRoot       string       `mapstructure:"merkle_root"`
	Logger           LoggerConfig `mapstructure:"logger"`
	HTTP             HTTPConfig   `mapstructure:"http"`
	Secret           SecretConfig `mapstructure:"secret"`
	Store            StoreConfig  `mapstructure:"store"`
	AccountWhiteList []string     `mapstructure:"account_white_list"`
}

func defaultConfig(v *viper.Viper) {
	v.SetDefault("chain_id", "Binance-Chain-Ganges")
	v.SetDefault("merkle_root", "0x0000000000000000000000000000000000000000000000000000000000000000")
	v.SetDefault("account_white_list", []string{})
}

type LoggerConfig struct {
	Level  string           `mapstructure:"level"`
	Format logger.LogFormat `mapstructure:"format"`
}

func defaultLoggerConfig(v *viper.Viper) {
	v.SetDefault("logger.level", "INFO")
	v.SetDefault("logger.format", logger.ConsoleFormat)
}

type HTTPConfig struct {
	Addr string `mapstructure:"addr"`
	Port uint16 `mapstructure:"port"`
}

func defaultHTTPConfig(v *viper.Viper) {
	v.SetDefault("http.addr", "0.0.0.0")
	v.SetDefault("http.port", 8080)
}

type SecretConfig struct {
	Type                   string                 `mapstructure:"type"`
	LocalSecretConfig      LocalSecretConfig      `mapstructure:"local_secret"`
	AWSSecretManagerConfig AWSSecretManagerConfig `mapstructure:"aws_secret_manager"`
}

type LocalSecretConfig struct {
	PrivateKey string `mapstructure:"private_key"`
}

type AWSSecretManagerConfig struct {
	Region     string `mapstructure:"region"`
	SecretName string `mapstructure:"secret_name"`
}

func defaultSecretConfig(v *viper.Viper) {
	v.SetDefault("secret.type", "local")
	v.SetDefault("secret.local_secret.private_key", "")
	v.SetDefault("secret.aws_secret_manager.region", "")
	v.SetDefault("secret.aws_secret_manager.secret_name", "")
}

type StoreConfig struct {
	Driver      string            `mapstructure:"driver"`
	MemoryStore MemoryStoreConfig `mapstructure:"memory_store"`
	SqlStore    SQLStoreConfig    `mapstructure:"sql_store"`
}

type MemoryStoreConfig struct {
	Accounts     string `mapstructure:"accounts"`
	MerkleProofs string `mapstructure:"merkle_proofs"`
}

type SQLStoreConfig struct {
	SQLDriver      string        `mapstructure:"sql_driver"`
	Host           string        `mapstructure:"host"`
	Port           uint          `mapstructure:"port"`
	DBName         string        `mapstructure:"dbname"`
	User           string        `mapstructure:"user"`
	Password       string        `mapstructure:"password"`
	ConnectTimeout string        `mapstructure:"connect_timeout"`
	ReadTimeout    string        `mapstructure:"read_timeout"`
	WriteTimeout   string        `mapstructure:"write_timeout"`
	DialTimeout    time.Duration `mapstructure:"dial_timeout"`
	MaxLifetime    time.Duration `mapstructure:"max_lifetime"`
	MaxIdleTime    time.Duration `mapstructure:"max_idletime"`
	MaxIdleConn    int           `mapstructure:"max_idle_conn"`
	MaxOpenConn    int           `mapstructure:"max_open_conn"`
	SSLMode        bool          `mapstructure:"ssl_mode"`
	LogLevel       int           `mapstructure:"log_level"`
}

func defaultStoreConfig(v *viper.Viper) {
	v.SetDefault("store.driver", "memory")

	// memory store
	v.SetDefault("store.memory_store.accounts", "./example/accounts.json")
	v.SetDefault("store.memory_store.merkle_proofs", "./example/merkle_proofs.json")

	// sql store
	v.SetDefault("store.sql_store.sql_driver", "mysql")
	v.SetDefault("store.sql_store.host", "localhost")
	v.SetDefault("store.sql_store.port", 3306)
	v.SetDefault("store.sql_store.dbname", "approver")
	v.SetDefault("store.sql_store.user", "root")
	v.SetDefault("store.sql_store.password", "")
	v.SetDefault("store.sql_store.connect_timeout", "10s")
	v.SetDefault("store.sql_store.read_timeout", "30s")
	v.SetDefault("store.sql_store.write_timeout", "30s")
	v.SetDefault("store.sql_store.dial_timeout", "10s")
	v.SetDefault("store.sql_store.max_idletime", "1h")
	v.SetDefault("store.sql_store.max_lifetime", "1h")
	v.SetDefault("store.sql_store.max_idle_conn", 2)
	v.SetDefault("store.sql_store.max_open_conn", 5)
	v.SetDefault("store.sql_store.log_level", gormLogger.Info)
	v.SetDefault("store.sql_store.ssl_mode", false)

}

func NewConfig(configPath string) (*Config, error) {
	var file *os.File
	file, err := os.Open(configPath)
	if len(configPath) > 0 && err != nil {
		return nil, err
	}

	v := viper.New()
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	/* default */
	defaultConfig(v)
	defaultLoggerConfig(v)
	defaultHTTPConfig(v)
	defaultSecretConfig(v)
	defaultStoreConfig(v)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.ReadConfig(file)

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return &config, err
	}

	return &config, nil
}
