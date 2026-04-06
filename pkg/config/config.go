package config

// Config with postgres connection params
type PgConfig struct {
	ConnStr         string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime int
	MaxConnIdleTime int
}
