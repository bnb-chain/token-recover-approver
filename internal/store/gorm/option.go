package gorm

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// An Option configures a gorm.DB
type Option interface {
	Apply(*gorm.DB) error
}

// OptionFunc is a function that configures a gorm.DB
type OptionFunc func(*gorm.DB) error

// Apply is a function that set value to gorm.DB
func (f OptionFunc) Apply(engine *gorm.DB) error {
	return f(engine)
}

func SetConnMaxIdleTime(maxIdleTime time.Duration) Option {
	return OptionFunc(func(engine *gorm.DB) error {
		sql, err := engine.DB()
		if err != nil {
			return err
		}
		sql.SetConnMaxIdleTime(maxIdleTime)
		return nil
	})
}

func SetConnMaxLifetime(maxlifetime time.Duration) Option {
	return OptionFunc(func(engine *gorm.DB) error {
		sql, err := engine.DB()
		if err != nil {
			return err
		}
		sql.SetConnMaxLifetime(maxlifetime)
		return nil
	})
}

func SetMaxIdleConns(maxIdleConns int) Option {
	return OptionFunc(func(engine *gorm.DB) error {
		sql, err := engine.DB()
		if err != nil {
			return err
		}
		sql.SetMaxIdleConns(maxIdleConns)
		return nil
	})
}

func SetMaxOpenConns(maxOpenConns int) Option {
	return OptionFunc(func(engine *gorm.DB) error {
		sql, err := engine.DB()
		if err != nil {
			return err
		}
		sql.SetMaxOpenConns(maxOpenConns)
		return nil
	})
}

func SetLogLevel(logLevel logger.LogLevel) Option {
	return OptionFunc(func(engine *gorm.DB) error {
		engine.Logger = engine.Logger.LogMode(logLevel)
		return nil
	})
}

func SetLogger(logger logger.Interface) Option {
	return OptionFunc(func(engine *gorm.DB) error {
		engine.Logger = logger
		return nil
	})
}
