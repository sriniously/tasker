package config

import (
	"fmt"
	"time"
)

type ObservabilityConfig struct {
	ServiceName  string             `koanf:"observability_service_name" validate:"required"`
	Environment  string             `koanf:"observability_environment" validate:"required"`
	Logging      LoggingConfig      `koanf:"observability_logging" validate:"required"`
	NewRelic     NewRelicConfig     `koanf:"observability_new_relic" validate:"required"`
	HealthChecks HealthChecksConfig `koanf:"observability_health_checks" validate:"required"`
}

type LoggingConfig struct {
	Level              string        `koanf:"observability_logging_level" validate:"required"`
	Format             string        `koanf:"observability_logging_format" validate:"required"`
	SlowQueryThreshold time.Duration `koanf:"observability_logging_slow_query_threshold"`
}

type NewRelicConfig struct {
	LicenseKey                string `koanf:"observability_new_relic_license_key" validate:"required"`
	AppLogForwardingEnabled   bool   `koanf:"observability_new_relic_app_log_forwarding_enabled"`
	DistributedTracingEnabled bool   `koanf:"observability_new_relic_distributed_tracing_enabled"`
	DebugLogging              bool   `koanf:"observability_new_relic_debug_logging"`
}

type HealthChecksConfig struct {
	Enabled  bool          `koanf:"observability_health_checks_enabled"`
	Interval time.Duration `koanf:"observability_health_checks_interval" validate:"min=1s"`
	Timeout  time.Duration `koanf:"observability_health_checks_timeout" validate:"min=1s"`
	Checks   []string      `koanf:"observability_health_checks_checks"`
}

func DefaultObservabilityConfig() *ObservabilityConfig {
	return &ObservabilityConfig{
		ServiceName: "tasker",
		Environment: "development",
		Logging: LoggingConfig{
			Level:              "info",
			Format:             "json",
			SlowQueryThreshold: 100 * time.Millisecond,
		},
		NewRelic: NewRelicConfig{
			LicenseKey:                "",
			AppLogForwardingEnabled:   true,
			DistributedTracingEnabled: true,
			DebugLogging:              false, // Disabled by default to avoid mixed log formats
		},
		HealthChecks: HealthChecksConfig{
			Enabled:  true,
			Interval: 30 * time.Second,
			Timeout:  5 * time.Second,
			Checks:   []string{"database", "redis"},
		},
	}
}

func (c *ObservabilityConfig) Validate() error {
	if c.ServiceName == "" {
		return fmt.Errorf("service_name is required")
	}

	// Validate log level
	validLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLevels[c.Logging.Level] {
		return fmt.Errorf("invalid logging level: %s (must be one of: debug, info, warn, error)", c.Logging.Level)
	}

	// Validate slow query threshold
	if c.Logging.SlowQueryThreshold < 0 {
		return fmt.Errorf("logging slow_query_threshold must be non-negative")
	}

	return nil
}

func (c *ObservabilityConfig) GetLogLevel() string {
	switch c.Environment {
	case "production":
		if c.Logging.Level == "" {
			return "info"
		}
	case "development":
		if c.Logging.Level == "" {
			return "debug"
		}
	}
	return c.Logging.Level
}

func (c *ObservabilityConfig) IsProduction() bool {
	return c.Environment == "production"
}
