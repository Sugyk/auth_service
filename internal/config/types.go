package config

import "time"

type AppConfig struct {
	DBCfg     *PgConfig
	HasherCfg *HasherConfig
	JWTConfig *JWTConfig
}

// Config with postgres connection params
type PgConfig struct {
	ConnStr         string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime int
	MaxConnIdleTime int
}

type HasherConfig struct {
	Cost int
}

type JWTConfig struct {
	TTL    time.Duration
	Secret string
}
