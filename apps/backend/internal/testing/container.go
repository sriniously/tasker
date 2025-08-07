package testing

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/sriniously/tasker/internal/config"
	"github.com/sriniously/tasker/internal/database"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDB struct {
	Pool      *pgxpool.Pool
	Container testcontainers.Container
	Config    *config.Config
}

// SetupTestDB creates a Postgres container and applies migrations
func SetupTestDB(t *testing.T) (*TestDB, func()) {
	t.Helper()

	ctx := context.Background()
	dbName := fmt.Sprintf("test_db_%s", uuid.New().String()[:8])
	dbUser := "testuser"
	dbPassword := "testpassword"

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       dbName,
			"POSTGRES_USER":     dbUser,
			"POSTGRES_PASSWORD": dbPassword,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(30 * time.Second),
	}

	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err, "failed to start postgres container")

	host, err := pgContainer.Host(ctx)
	require.NoError(t, err, "failed to get container host")

	mappedPort, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err, "failed to get mapped port")
	port := mappedPort.Int()

	// Make sure the test cleans up the container
	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	})

	// Create configuration
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:            host,
			Port:            port,
			User:            dbUser,
			Password:        dbPassword,
			Name:            dbName,
			SSLMode:         "disable",
			MaxOpenConns:    25,
			MaxIdleConns:    25,
			ConnMaxLifetime: 300,
			ConnMaxIdleTime: 300,
		},
		Primary: config.Primary{
			Env: "test",
		},
		Server: config.ServerConfig{
			Port:               "8080",
			ReadTimeout:        30,
			WriteTimeout:       30,
			IdleTimeout:        30,
			CORSAllowedOrigins: []string{"*"},
		},
		Integration: config.IntegrationConfig{
			ResendAPIKey: "test-key",
		},
		Redis: config.RedisConfig{
			Address: "localhost:6379",
		},
		Auth: config.AuthConfig{
			SecretKey: "test-secret",
		},
	}

	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	var db *database.Database
	var lastErr error
	for i := 0; i < 5; i++ {
		// Sleep before first attempt too to give PostgreSQL time to initialize
		time.Sleep(2 * time.Second)

		db, lastErr = database.New(cfg, &logger, nil)
		if lastErr == nil {
			// Try a ping to verify the connection
			if err := db.Pool.Ping(ctx); err == nil {
				break
			} else {
				lastErr = err
				logger.Warn().Err(err).Msg("Failed to ping database, will retry")
				db.Pool.Close() // Close the failed connection
			}
		} else {
			logger.Warn().Err(lastErr).Msgf("Failed to connect to database (attempt %d/5)", i+1)
		}
	}
	require.NoError(t, lastErr, "failed to connect to database after multiple attempts")

	// Apply migrations
	err = database.Migrate(ctx, &logger, cfg)
	require.NoError(t, err, "failed to apply database migrations")

	testDB := &TestDB{
		Pool:      db.Pool,
		Container: pgContainer,
		Config:    cfg,
	}

	// Return cleanup function that just closes the pool (container is managed by t.Cleanup)
	cleanup := func() {
		if db.Pool != nil {
			db.Pool.Close()
		}
	}

	return testDB, cleanup
}

// CleanupTestDB closes the database connection and terminates the container
func (db *TestDB) CleanupTestDB(ctx context.Context, logger *zerolog.Logger) error {
	logger.Info().Msg("cleaning up test database")

	if db.Pool != nil {
		db.Pool.Close()
	}

	if db.Container != nil {
		if err := db.Container.Terminate(ctx); err != nil {
			return fmt.Errorf("failed to terminate container: %w", err)
		}
	}

	return nil
}
