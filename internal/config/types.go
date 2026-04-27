package config

type AppConfig struct {
	DBCfg     *PgConfig
	HasherCfg *HasherConfig
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
