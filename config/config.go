package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	EnvDbHost        = "DB_HOST"
	EnvDbName        = "DB_NAME"
	EnvDbPort        = "DB_PORT"
	EnvDbUser        = "DB_USER"
	EnvDbPswd        = "DB_PASSWORD"
	EnvRedisAddr     = "REDIS_ADDR"
	EnvRedisPassword = "REDIS_PASSWORD"
	EnvRedisDB       = "REDIS_DB"
	EnvCacheTTL      = "CACHE_TTL_SECONDS"
)

type Config struct {
	DbHost        string
	DbName        string
	DbPort        string
	DbUser        string
	DbPswd        string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	CacheTTL      time.Duration
}

func GetFromEnv() *Config {
	conf := &Config{}
	flag.StringVar(&conf.DbHost, EnvDbHost, os.Getenv(EnvDbHost), "database host")
	flag.StringVar(&conf.DbName, EnvDbName, os.Getenv(EnvDbName), "database name")
	flag.StringVar(&conf.DbPort, EnvDbPort, os.Getenv(EnvDbPort), "database port")
	flag.StringVar(&conf.DbUser, EnvDbUser, os.Getenv(EnvDbUser), "database user name")
	flag.StringVar(&conf.DbPswd, EnvDbPswd, os.Getenv(EnvDbPswd), "database user password")
	flag.StringVar(&conf.RedisAddr, EnvRedisAddr, getEnvOrDefault(EnvRedisAddr, "localhost:6379"), "redis address")
	flag.StringVar(&conf.RedisPassword, EnvRedisPassword, os.Getenv(EnvRedisPassword), "redis password")
	flag.Parse()

	conf.RedisDB = getEnvAsInt(EnvRedisDB, 0)
	conf.CacheTTL = time.Duration(getEnvAsInt(EnvCacheTTL, 300)) * time.Second

	return conf
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func (cfg *Config) ConnString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DbHost,
		cfg.DbUser,
		cfg.DbPswd,
		cfg.DbName,
		cfg.DbPort,
	)
}
