package config

import "time"

type AppConfig struct {
	DBCfg       *PgConfig
	HasherCfg   *HasherConfig
	JWTConfig   *JWTConfig
	GRPCConfig  *GRPCConfig
	RedisCfg    *RedisConfig
	ThrottleCfg *LoginThrottleConfig
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

type GRPCConfig struct {
	Addr string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// Config for the login brute-force throttle: after MaxAttempts failed logins
// within BlockDuration, further attempts are blocked for BlockDuration.
type LoginThrottleConfig struct {
	MaxAttempts   int
	BlockDuration time.Duration
}
