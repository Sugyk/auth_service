package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// LoadConfig загружает конфигурацию
// Сначала пытается прочитать yaml-файл, если его нет — читает из переменных окружения
func LoadConfig() (*AppConfig, error) {
	// Поддержка .env файла (удобно для разработки)
	_ = godotenv.Load()

	v := viper.New()

	// Настройка для YAML
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("./app/config")
	v.AddConfigPath(filepath.Join(".", "config"))

	// Настройка маппинга ключей: BindEnv с двумя аргументами берёт имя
	// переменной окружения как есть, поэтому префикс APP_ указан явно.
	v.BindEnv("pg.connstr", "APP_PG_CONNSTR")
	v.BindEnv("pg.maxconns", "APP_PG_MAX_CONNS")
	v.BindEnv("pg.minconns", "APP_PG_MIN_CONNS")
	v.BindEnv("pg.maxconnlifetime", "APP_PG_MAX_CONN_LIFETIME")
	v.BindEnv("pg.maxconnidletime", "APP_PG_MAX_CONN_IDLE_TIME")

	v.BindEnv("hasher.cost", "APP_HASHER_COST")

	v.BindEnv("jwt.ttl", "APP_JWT_TTL")

	v.BindEnv("grpc.addr", "APP_GRPC_ADDR")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		fmt.Println("Config file not found, using environment variables only")
	} else {
		fmt.Printf("Config loaded from file: %s\n", v.ConfigFileUsed())
	}

	pgConfig := &PgConfig{
		ConnStr:         v.GetString("pg.connstr"),
		MaxConns:        v.GetInt32("pg.maxconns"),
		MinConns:        v.GetInt32("pg.minconns"),
		MaxConnLifetime: v.GetInt("pg.maxconnlifetime"),
		MaxConnIdleTime: v.GetInt("pg.maxconnidletime"),
	}

	hasherConfig := &HasherConfig{
		Cost: v.GetInt("hasher.cost"),
	}

	jwtConfig := &JWTConfig{
		Secret: os.Getenv("JWT_SECRET"),
		TTL:    v.GetDuration("jwt.ttl"),
	}

	grpcConfig := &GRPCConfig{
		Addr: v.GetString("grpc.addr"),
	}

	if pgConfig.ConnStr == "" {
		return nil, fmt.Errorf("PG_CONNSTR (or pg.connstr in yaml) is required")
	}

	setDefaults(pgConfig, hasherConfig, grpcConfig)

	appConfig := &AppConfig{
		DBCfg:      pgConfig,
		HasherCfg:  hasherConfig,
		JWTConfig:  jwtConfig,
		GRPCConfig: grpcConfig,
	}

	return appConfig, nil
}

func setDefaults(pg *PgConfig, hasher *HasherConfig, grpcCfg *GRPCConfig) {
	if pg.MaxConns <= 0 {
		pg.MaxConns = 25
	}
	if pg.MinConns <= 0 {
		pg.MinConns = 2
	}
	if pg.MaxConnLifetime <= 0 {
		pg.MaxConnLifetime = 30 * 60 // 30 минут
	}
	if pg.MaxConnIdleTime <= 0 {
		pg.MaxConnIdleTime = 5 * 60 // 5 минут
	}
	if hasher.Cost <= 0 {
		hasher.Cost = 12
	}
	if grpcCfg.Addr == "" {
		grpcCfg.Addr = ":50051"
	}
}
