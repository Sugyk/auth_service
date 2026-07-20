package integration_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Sugyk/auth_service/internal/api/http/handlers"
	"github.com/Sugyk/auth_service/internal/config"
	"github.com/Sugyk/auth_service/internal/pkg/hasher"
	"github.com/Sugyk/auth_service/internal/pkg/jwt_manager"
	"github.com/Sugyk/auth_service/internal/pkg/ratelimiter"
	"github.com/Sugyk/auth_service/internal/repository"
	"github.com/Sugyk/auth_service/internal/service"
	"github.com/Sugyk/auth_service/pkg/logger"
	"github.com/Sugyk/auth_service/pkg/postgres"
	"github.com/alicebob/miniredis/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type IntegrationSuite struct {
	suite.Suite

	cfg *config.AppConfig

	handler *handlers.Handler

	db *pgxpool.Pool
	tx pgx.Tx

	mr *miniredis.Miniredis
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}

func (s *IntegrationSuite) PerformRequest(
	method string,
	url string,
	body any,
	handler http.HandlerFunc,
) *httptest.ResponseRecorder {

	var payload []byte

	if body != nil {
		var err error

		payload, err = json.Marshal(body)
		s.Require().NoError(err)
	}

	req := httptest.NewRequest(
		method,
		url,
		bytes.NewReader(payload),
	)

	rr := httptest.NewRecorder()

	handler(rr, req)

	return rr
}

func (s *IntegrationSuite) SetupSuite() {
	cfg, err := config.LoadConfig()
	s.Require().NoError(err)

	s.cfg = cfg

	log := logger.NewNoop()

	pgProvider := postgres.NewProvider(
		log,
		cfg.DBCfg.ConnStr,
		cfg.DBCfg.MaxConns,
		cfg.DBCfg.MinConns,
		cfg.DBCfg.MaxConnLifetime,
		cfg.DBCfg.MaxConnIdleTime,
	)

	for range 10 {
		err = pgProvider.Open(s.T().Context())
		if err != nil {
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}

	s.Require().NoError(err)

	s.db = pgProvider.DB()

	mr, err := miniredis.Run()
	s.Require().NoError(err)

	s.mr = mr
}

func (s *IntegrationSuite) TearDownSuite() {
	if s.mr != nil {
		s.mr.Close()
	}
}

func (s *IntegrationSuite) SetupTest() {
	tx, err := s.db.BeginTx(s.T().Context(), pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	s.Require().NoError(err)

	s.tx = tx

	txManager := postgres.NewTestTxManager(tx)

	jwtManager, err := jwt_manager.NewJWTManager(
		[]byte(s.cfg.JWTConfig.Secret),
		s.cfg.JWTConfig.TTL,
	)
	s.Require().NoError(err)

	repo := repository.NewRepository(s.db)

	passwordHasher := hasher.NewPasswordHasher(
		s.cfg.HasherCfg.Cost,
	)

	s.mr.FlushAll()
	redisClient := redis.NewClient(&redis.Options{Addr: s.mr.Addr()})
	throttler := ratelimiter.New(
		redisClient,
		s.cfg.ThrottleCfg.MaxAttempts,
		s.cfg.ThrottleCfg.BlockDuration,
	)

	svc := service.NewService(
		repo,
		txManager,
		passwordHasher,
		jwtManager,
		throttler,
	)

	s.handler = handlers.NewHandler(
		svc,
		logger.NewNoop(),
	)
}

func (s *IntegrationSuite) TearDownTest() {
	if s.tx != nil {
		_ = s.tx.Rollback(s.T().Context())
	}
}
