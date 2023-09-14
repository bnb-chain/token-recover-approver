package config

import (
	"os"
	"strings"

	"github.com/bnb-chain/airdrop-service/pkg/logger"

	"github.com/spf13/viper"
)

type Config struct {
	ChainID string       `mapstructure:"chain_id"`
	Logger  LoggerConfig `mapstructure:"logger"`
	HTTP    HTTPConfig   `mapstructure:"http"`
	Secret  SecretConfig `mapstructure:"secret"`
	Store   StoreConfig  `mapstructure:"store"`
}

func defaultConfig(v *viper.Viper) {
	v.SetDefault("chain_id", "Binance-Chain-Ganges")
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
}

type MemoryStoreConfig struct {
	StateRoot    string `mapstructure:"state_root"`
	Assets       string `mapstructure:"assets"`
	Accounts     string `mapstructure:"accounts"`
	MerkleProofs string `mapstructure:"merkle_proofs"`
}

func defaultStoreConfig(v *viper.Viper) {
	v.SetDefault("store.driver", "memory")
	v.SetDefault("store.memory_store.state_root", "./example/state_root.json")
	v.SetDefault("store.memory_store.assets", "./example/assets.json")
	v.SetDefault("store.memory_store.accounts", "./example/accounts.json")
	v.SetDefault("store.memory_store.merkle_proofs", "./example/merkle_proofs.json")
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
