package config

import (
	"fmt"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// LoadConfig загружает конфигурацию
// Сначала пытается прочитать yaml-файл, если его нет — читает из переменных окружения
func LoadConfig() (*PgConfig, *HasherConfig, error) {
	// Поддержка .env файла (удобно для разработки)
	_ = godotenv.Load()

	v := viper.New()

	// Настройка для YAML
	v.SetConfigName("config") // имя файла без расширения
	v.SetConfigType("yaml")
	v.AddConfigPath(".")         // текущая директория
	v.AddConfigPath("./config")  // папка config
	v.AddConfigPath("./configs") // часто используют configs (мн.ч.)
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

	// Читаем конфиг файл (если он существует)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Это другая ошибка (неправильный yaml и т.д.) — возвращаем
			return nil, nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Файл не найден — это нормально, будем читать только из ENV
		fmt.Println("Config file not found, using environment variables only")
	} else {
		fmt.Printf("Config loaded from file: %s\n", v.ConfigFileUsed())
	}

	// Заполняем структуры
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

	// Валидация обязательных полей
	if pgConfig.ConnStr == "" {
		return nil, nil, fmt.Errorf("PG_CONNSTR (or pg.connstr in yaml) is required")
	}

	// Установка разумных значений по умолчанию, если не указаны
	setDefaults(pgConfig, hasherConfig)

	return pgConfig, hasherConfig, nil
}

// setDefaults устанавливает разумные значения по умолчанию
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
