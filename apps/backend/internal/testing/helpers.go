package testing

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/sriniously/tasker/internal/server"
	"github.com/stretchr/testify/require"
)

// SetupTest prepares a test environment with a database and server
func SetupTest(t *testing.T) (*TestDB, *server.Server, func()) {
	t.Helper()

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()

	testDB, dbCleanup := SetupTestDB(t)

	testServer := CreateTestServer(&logger, testDB)

	cleanup := func() {
		if testDB.Pool != nil {
			testDB.Pool.Close()
		}

		dbCleanup()
	}

	return testDB, testServer, cleanup
}

// MustMarshalJSON marshals an object to JSON or fails the test
func MustMarshalJSON(t *testing.T, v interface{}) []byte {
	t.Helper()

	jsonBytes, err := json.Marshal(v)
	require.NoError(t, err, "failed to marshal to JSON")

	return jsonBytes
}

// ProjectRoot returns the absolute path to the project root
func ProjectRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	require.NoError(t, err, "failed to get working directory")

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			t.Fatal("could not find project root (go.mod)")
			return ""
		}

		dir = parentDir
	}
}

// Ptr returns a pointer to the given value
// Useful for creating pointers to values for optional fields
func Ptr[T any](v T) *T {
	return &v
}
