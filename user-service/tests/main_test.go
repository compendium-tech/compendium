package tests

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	"github.com/compendium-tech/compendium/common/pkg/pg"
	"github.com/compendium-tech/compendium/common/pkg/redis"
	"github.com/compendium-tech/compendium/common/pkg/validate"

	"github.com/compendium-tech/compendium/user-service/internal/app"
	"github.com/compendium-tech/compendium/user-service/internal/config"
	"github.com/compendium-tech/compendium/user-service/internal/email"
	"github.com/compendium-tech/compendium/user-service/internal/geoip"
	"github.com/compendium-tech/compendium/user-service/internal/hash"
	"github.com/compendium-tech/compendium/user-service/internal/ua"
)

type APITestSuite struct {
	suite.Suite
	app.GinAppDependencies
	ctx                     context.Context
	app                     *gin.Engine
	mockEmailMessageBuilder *email.MockMessageBuilder
	mockEmailSender         *email.MockSender
	mockGeoIP               *geoip.MockGeoIP
	mockUserAgentParser     *ua.MockUserAgentParser
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
	s.initDeps()

	if err := pg.RunUpPgMigrations(s.ctx, s.PgDB, s.getPgMigrationsDir()); err != nil {
		s.FailNow("Failed to run up postgres migrations", err)
	}
}

func (s *APITestSuite) TearDownTest() {
	if err := pg.RunDownPgMigrations(s.ctx, s.PgDB, s.getPgMigrationsDir()); err != nil {
		s.FailNow("Failed to run down postgres migrations", err)
	}

	if s.PgDB != nil {
		s.T().Log("Closing PostgreSQL connection...")
		if err := s.PgDB.Close(); err != nil {
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

func (s *APITestSuite) initDeps() {
	validate.InitValidator()

	ctx := context.Background()
	err := godotenv.Load(".env")
	if err != nil {
		s.T().Logf("Failed to load .env file, using environmental variables instead: %v\n", err)
	}

	cfg := config.LoadGinAppConfig()

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
	s.mockEmailSender = new(email.MockSender)
	s.mockEmailMessageBuilder = new(email.MockMessageBuilder)
	s.mockGeoIP = new(geoip.MockGeoIP)
	s.mockUserAgentParser = new(ua.MockUserAgentParser)

	s.GinAppDependencies = app.GinAppDependencies{
		PgDB:                pgDB,
		RedisClient:         redisClient,
		Config:              cfg,
		TokenManager:        tokenManager,
		EmailSender:         s.mockEmailSender,
		EmailMessageBuilder: s.mockEmailMessageBuilder,
		GeoIP:               s.mockGeoIP,
		UserAgentParser:     s.mockUserAgentParser,
		PasswordHasher:      hash.NewBcryptPasswordHasher(bcrypt.DefaultCost),
	}

	s.app = app.NewGinApp(s.GinAppDependencies)
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
