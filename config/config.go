package config

import (
	"context"
	"time"

	"github.com/sethvargo/go-envconfig"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type ContextKey string

const ContextUserIDKey ContextKey = "userID"

const (
	mongoDBPrefix = "MONGO_DB_"
	serverPrefix  = "SERVER_"
)

type DBConfig struct {
	Host     string `env:"HOST, required"`
	Port     string `env:"PORT, required"`
	User     string `env:"USER, required"`
	Password string `env:"PASSWORD, required"`
	DBName   string `env:"DB_NAME, required"`
}

type MongoDBConfig struct {
	URI        string `env:"URI, required"`
	Password   string `env:"PASSWORD"`
	User       string `env:"USER"`
	Database   string `env:"DATABASE, required"`
	Collection string `env:"COLLECTION, required"`
}

type CacheConfig struct {
	Addr     string `env:"ADDR, required"`
	Password string `env:"PASSWORD, required"`
	DB       int    `env:"DB, required"`
}

type ServerConfig struct {
	ReadTimeout               time.Duration `env:"READ_TIMEOUT,default=5s"`
	WriteTimeout              time.Duration `env:"WRITE_TIMEOUT,default=10s"`
	IdleTimeout               time.Duration `env:"IDLE_TIMEOUT,default=15s"`
	GracefulShutdownTimeout   time.Duration `env:"GRACEFUL_SHUTDOWN_TIMEOUT,default=5s"`
	PreflightCacheDurationSec int           `env:"PREFLIGHT_CACHE_DURATION_SEC,default=300"`
}

func NewDBConfig(ctx context.Context, dbPrefix string) (*DBConfig, error) {
	conf := &DBConfig{}
	pl := envconfig.PrefixLookuper(dbPrefix, envconfig.OsLookuper())
	if err := envconfig.ProcessWith(ctx, conf, pl); err != nil {
		log.Error("Failed to load database config", log.Ferror(err))
		return nil, err
	}
	return conf, nil
}

func NewMongoDBConfig(ctx context.Context) (*MongoDBConfig, error) {
	conf := &MongoDBConfig{}
	pl := envconfig.PrefixLookuper(mongoDBPrefix, envconfig.OsLookuper())
	if err := envconfig.ProcessWith(ctx, conf, pl); err != nil {
		log.Error("Failed to load MongoDB config", log.Ferror(err))
		return nil, err
	}

	return conf, nil
}

func NewCacheConfig(ctx context.Context, cachePrefix string) (*CacheConfig, error) {
	conf := &CacheConfig{}
	pl := envconfig.PrefixLookuper(cachePrefix, envconfig.OsLookuper())
	if err := envconfig.ProcessWith(ctx, conf, pl); err != nil {
		log.Error("Failed to load cache config", log.Ferror(err))
		return nil, err
	}
	return conf, nil
}

func NewServerConfig(ctx context.Context) (*ServerConfig, error) {
	conf := &ServerConfig{}
	pl := envconfig.PrefixLookuper(serverPrefix, envconfig.OsLookuper())
	if err := envconfig.ProcessWith(ctx, conf, pl); err != nil {
		log.Error("Failed to load server config", log.Ferror(err))
		return nil, err
	}
	return conf, nil
}
