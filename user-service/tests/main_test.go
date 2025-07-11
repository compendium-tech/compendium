package tests

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/seacite-tech/compendium/common/pkg/pg"
	"github.com/seacite-tech/compendium/common/pkg/redis"
	"github.com/seacite-tech/compendium/user-service/internal/app"
	"github.com/seacite-tech/compendium/user-service/internal/config"
	"github.com/seacite-tech/compendium/user-service/internal/email"
	"github.com/seacite-tech/compendium/user-service/internal/validate"
	"github.com/seacite-tech/compendium/user-service/pkg/auth"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	app.Dependencies
	ctx             context.Context
	app             *gin.Engine
	mockEmailSender *email.MockEmailSender
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping API suite in short mode.")
	}

	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupTest() {
	s.initializeDeps()

	if err := pg.RunUpPgMigrations(s.ctx, s.PgDb, s.getPgMigrationsDir()); err != nil {
		s.FailNow("Failed to run up postgres migrations", err)
	}
}

func (s *APITestSuite) TearDownTest() {
	if err := pg.RunDownPgMigrations(s.ctx, s.PgDb, s.getPgMigrationsDir()); err != nil {
		s.FailNow("Failed to run down postgres migrations", err)
	}

	if s.PgDb != nil {
		s.T().Log("Closing PostgreSQL connection...")
		if err := s.PgDb.Close(); err != nil {
			s.T().Errorf("Error closing PostgreSQL connection: %v", err)
		}
	}

	s.T().Log("Flushing Redis database and closing client connection...")
	if err := s.RedisClient.FlushDB(s.ctx).Err(); err != nil {
		s.T().Errorf("Error flushing Redis database: %v", err)
	}

	if s.RedisClient != nil {
		s.T().Log("Closing Redis client connection...")
		if err := s.RedisClient.Close(); err != nil {
			s.T().Errorf("Error closing Redis client connection: %v", err)
		}
	}
}

func (s *APITestSuite) initializeDeps() {
	validate.InitializeValidator()

	ctx := context.Background()
	err := godotenv.Load(".env")
	if err != nil {
		s.T().Logf("Failed to load .env file, using environmental variables instead: %v\n", err)
	}

	cfg := config.LoadAppConfig()

	tokenManager, err := auth.NewJwtBasedTokenManager(cfg.JwtSingingKey)
	if err != nil {
		s.FailNow("Failed to initialize token manager, cause: %s", err)
		return
	}

	pgDB, err := pg.NewPgClient(ctx, cfg.PgHost, cfg.PgPort, cfg.PgUsername, cfg.PgPassword, cfg.PgDatabaseName)
	if err != nil {
		s.FailNow("Failed to connect to PostgreSQL, cause: %s", err)
		return
	}

	redisClient, err := redis.NewRedisClient(ctx, cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		s.FailNow("Failed to connect to Redis, cause: %s", err)
		return
	}

	s.ctx = ctx
	s.mockEmailSender = new(email.MockEmailSender)
	s.Dependencies = app.Dependencies{
		PgDb:         pgDB,
		RedisClient:  redisClient,
		Config:       cfg,
		TokenManager: tokenManager,
		EmailSender:  s.mockEmailSender,
	}
	s.app = app.NewApp(s.Dependencies)
}

func (s *APITestSuite) getPgMigrationsDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		s.FailNow("failed to get current file path")
	}

	currentDir := filepath.Dir(filename)
	migrationsDir := filepath.Join(currentDir, "..", "migrations")

	s.T().Logf("Migrations directory: %s", migrationsDir)

	return migrationsDir
}
