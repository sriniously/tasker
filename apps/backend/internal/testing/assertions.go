package testing

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertTimestampsValid checks that created_at and updated_at fields are set
func AssertTimestampsValid(t *testing.T, obj interface{}) {
	t.Helper()

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	createdField := val.FieldByName("CreatedAt")
	if createdField.IsValid() {
		createdAt, ok := createdField.Interface().(time.Time)
		require.True(t, ok, "CreatedAt is not a time.Time")
		assert.False(t, createdAt.IsZero(), "CreatedAt should not be zero")
	}

	updatedField := val.FieldByName("UpdatedAt")
	if updatedField.IsValid() {
		updatedAt, ok := updatedField.Interface().(time.Time)
		require.True(t, ok, "UpdatedAt is not a time.Time")
		assert.False(t, updatedAt.IsZero(), "UpdatedAt should not be zero")
	}
}

// AssertValidUUID checks that the UUID is valid and not nil
func AssertValidUUID(t *testing.T, id uuid.UUID, message ...string) {
	t.Helper()

	msg := "UUID should not be nil"
	if len(message) > 0 {
		msg = message[0]
	}

	assert.NotEqual(t, uuid.Nil, id, msg)
}

// AssertEqualExceptTime asserts that two objects are equal, ignoring time fields
func AssertEqualExceptTime(t *testing.T, expected, actual interface{}) {
	t.Helper()

	expectedVal := reflect.ValueOf(expected)
	if expectedVal.Kind() == reflect.Ptr {
		expectedVal = expectedVal.Elem()
	}

	actualVal := reflect.ValueOf(actual)
	if actualVal.Kind() == reflect.Ptr {
		actualVal = actualVal.Elem()
	}

	// Ensure same type
	require.Equal(t, expectedVal.Type(), actualVal.Type(), "objects are not the same type")

	// Check fields
	for i := 0; i < expectedVal.NumField(); i++ {
		field := expectedVal.Type().Field(i)

		// Skip time fields
		if field.Type == reflect.TypeOf(time.Time{}) ||
			field.Type == reflect.TypeOf(&time.Time{}) {
			continue
		}

		expectedField := expectedVal.Field(i)
		actualField := actualVal.Field(i)

		assert.Equal(
			t,
			expectedField.Interface(),
			actualField.Interface(),
			fmt.Sprintf("field %s should be equal", field.Name),
		)
	}
}

// AssertStringContains checks if a string contains all specified substrings
func AssertStringContains(t *testing.T, s string, substrings ...string) {
	t.Helper()

	for _, sub := range substrings {
		assert.True(
			t,
			strings.Contains(s, sub),
			fmt.Sprintf("expected string to contain '%s', but it didn't: %s", sub, s),
		)
	}
}
