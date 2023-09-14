package config

import (
	"os"
	"strings"

	"github.com/bnb-chain/airdrop-service/pkg/logger"

	"github.com/spf13/viper"
)

type Config struct {
	Logger LoggerConfig `mapstructure:"logger"`
	HTTP   HTTPConfig   `mapstructure:"http"`
	Secret SecretConfig `mapstructure:"secret"`
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

func NewConfig(configPath string) (*Config, error) {
	var file *os.File
	file, _ = os.Open(configPath)

	v := viper.New()
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	/* default */
	defaultLoggerConfig(v)
	defaultHTTPConfig(v)
	defaultSecretConfig(v)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.ReadConfig(file)

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return &config, err
	}

	return &config, nil
}
