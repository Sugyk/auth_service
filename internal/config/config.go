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

	// Автоматическое приведение имён переменных окружения
	v.AutomaticEnv()
	v.SetEnvPrefix("APP") // префикс: APP_PG_CONNSTR, APP_HASHER_COST и т.д.

	// Настройка маппинга ключей (для удобства)
	v.BindEnv("pg.connstr", "PG_CONNSTR")
	v.BindEnv("pg.maxconns", "PG_MAX_CONNS")
	v.BindEnv("pg.minconns", "PG_MIN_CONNS")
	v.BindEnv("pg.maxconnlifetime", "PG_MAX_CONN_LIFETIME")
	v.BindEnv("pg.maxconnidletime", "PG_MAX_CONN_IDLE_TIME")

	v.BindEnv("hasher.cost", "HASHER_COST")

	v.BindEnv("jwt.ttl", "JWT_TTL")

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

	if pgConfig.ConnStr == "" {
		return nil, fmt.Errorf("PG_CONNSTR (or pg.connstr in yaml) is required")
	}

	setDefaults(pgConfig, hasherConfig)

	appConfig := &AppConfig{
		DBCfg:     pgConfig,
		HasherCfg: hasherConfig,
		JWTConfig: jwtConfig,
	}

	return appConfig, nil
}

func setDefaults(pg *PgConfig, hasher *HasherConfig) {
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
}
