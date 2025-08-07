package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"github.com/sriniously/tasker/internal/logger"
	"github.com/sriniously/tasker/internal/server"
)

const (
	UserIDKey   = "user_id"
	UserRoleKey = "user_role"
	LoggerKey   = "logger"
)

type ContextEnhancer struct {
	server *server.Server
}

func NewContextEnhancer(s *server.Server) *ContextEnhancer {
	return &ContextEnhancer{server: s}
}

func (ce *ContextEnhancer) EnhanceContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract request ID
			requestID := GetRequestID(c)

			// Create enhanced logger with request context
			contextLogger := ce.server.Logger.With().
				Str("request_id", requestID).
				Str("method", c.Request().Method).
				Str("path", c.Path()).
				Str("ip", c.RealIP()).
				Logger()

			// Add trace context if available
			if txn := newrelic.FromContext(c.Request().Context()); txn != nil {
				contextLogger = logger.WithTraceContext(contextLogger, txn)
			}

			// Extract user information from JWT token or session
			if userID := ce.extractUserID(c); userID != "" {
				contextLogger = contextLogger.With().Str("user_id", userID).Logger()
			}

			if userRole := ce.extractUserRole(c); userRole != "" {
				contextLogger = contextLogger.With().Str("user_role", userRole).Logger()
			}

			// Store the enhanced logger in context
			c.Set(LoggerKey, &contextLogger)

			// Create a new context with the logger
			ctx := context.WithValue(c.Request().Context(), LoggerKey, &contextLogger)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func (ce *ContextEnhancer) extractUserID(c echo.Context) string {
	// Check if user_id was already set by auth middleware (Clerk)
	if userID, ok := c.Get("user_id").(string); ok && userID != "" {
		return userID
	}
	return ""
}

func (ce *ContextEnhancer) extractUserRole(c echo.Context) string {
	// Check if user_role was set by auth middleware (Clerk)
	if userRole, ok := c.Get("user_role").(string); ok && userRole != "" {
		return userRole
	}
	return ""
}

func GetUserID(c echo.Context) string {
	if userID, ok := c.Get(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

func GetLogger(c echo.Context) *zerolog.Logger {
	if logger, ok := c.Get(LoggerKey).(*zerolog.Logger); ok {
		return logger
	}
	// Fallback to a basic logger if not found
	logger := zerolog.Nop()
	return &logger
}
