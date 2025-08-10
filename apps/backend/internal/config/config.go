package config

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	Primary       Primary              `koanf:"primary" validate:"required"`
	Server        ServerConfig         `koanf:"server" validate:"required"`
	Database      DatabaseConfig       `koanf:"database" validate:"required"`
	Auth          AuthConfig           `koanf:"auth" validate:"required"`
	Redis         RedisConfig          `koanf:"redis" validate:"required"`
	Integration   IntegrationConfig    `koanf:"integration" validate:"required"`
	Observability *ObservabilityConfig `koanf:"observability"`
	AWS           AWSConfig            `koanf:"aws" validate:"required"`
	Cron          *CronConfig          `koanf:"cron"`
}

type Primary struct {
	Env string `koanf:"primary_env" validate:"required"`
}

type ServerConfig struct {
	Port               string   `koanf:"server_port" validate:"required"`
	ReadTimeout        int      `koanf:"server_read_timeout" validate:"required"`
	WriteTimeout       int      `koanf:"server_write_timeout" validate:"required"`
	IdleTimeout        int      `koanf:"server_idle_timeout" validate:"required"`
	CORSAllowedOrigins []string `koanf:"server_cors_allowed_origins" validate:"required"`
}

type DatabaseConfig struct {
	Host            string `koanf:"database_host" validate:"required"`
	Port            int    `koanf:"database_port" validate:"required"`
	User            string `koanf:"database_user" validate:"required"`
	Password        string `koanf:"database_password"`
	Name            string `koanf:"database_name" validate:"required"`
	SSLMode         string `koanf:"database_ssl_mode" validate:"required"`
	MaxOpenConns    int    `koanf:"database_max_open_conns" validate:"required"`
	MaxIdleConns    int    `koanf:"database_max_idle_conns" validate:"required"`
	ConnMaxLifetime int    `koanf:"database_conn_max_lifetime" validate:"required"`
	ConnMaxIdleTime int    `koanf:"database_conn_max_idle_time" validate:"required"`
}
type RedisConfig struct {
	Address string `koanf:"redis_address" validate:"required"`
}

type IntegrationConfig struct {
	ResendAPIKey string `koanf:"integration_resend_api_key" validate:"required"`
}

type AuthConfig struct {
	SecretKey string `koanf:"auth_secret_key" validate:"required"`
}

type AWSConfig struct {
	Region          string `koanf:"aws_region" validate:"required"`
	AccessKeyID     string `koanf:"aws_access_key_id" validate:"required"`
	SecretAccessKey string `koanf:"aws_secret_access_key" validate:"required"`
	UploadBucket    string `koanf:"aws_upload_bucket" validate:"required"`
	EndpointURL     string `koanf:"aws_endpoint_url"`
}

type CronConfig struct {
	ArchiveDaysThreshold        int `koanf:"cron_archive_days_threshold"`
	BatchSize                   int `koanf:"cron_batch_size"`
	ReminderHours               int `koanf:"cron_reminder_hours"`
	MaxTodosPerUserNotification int `koanf:"cron_max_todos_per_user_notification"`
}

func DefaultCronConfig() *CronConfig {
	return &CronConfig{
		ArchiveDaysThreshold:        30,
		BatchSize:                   100,
		ReminderHours:               24,
		MaxTodosPerUserNotification: 10,
	}
}

func LoadConfig() (*Config, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	k := koanf.New(".")

	err := k.Load(env.Provider("TASKER_", "", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "TASKER_"))
	}), nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not load initial env variables")
	}

	mainConfig := &Config{}

	configTargets := []struct {
		target interface{}
		name   string
	}{
		{&mainConfig.AWS, "aws config"},
		{&mainConfig.Database, "database config"},
		{&mainConfig.Primary, "primary config"},
		{&mainConfig.Auth, "auth config"},
		{&mainConfig.Redis, "redis config"},
		{&mainConfig.Integration, "integration config"},
		{&mainConfig.Cron, "cron config"},
		{&mainConfig.Server, "server config"},
	}

	for _, target := range configTargets {
		if err := k.Unmarshal("", target.target); err != nil {
			logger.Fatal().Err(err).Msgf("could not unmarshal %s", target.name)
		}
	}

	if err := k.Unmarshal("", &mainConfig.Observability); err != nil {
		logger.Fatal().Err(err).Msg("could not unmarshal observability config")
	}
	if mainConfig.Observability != nil {
		if err := k.Unmarshal("", &mainConfig.Observability.Logging); err != nil {
			logger.Fatal().Err(err).Msg("could not unmarshal observability logging config")
		}
		if err := k.Unmarshal("", &mainConfig.Observability.NewRelic); err != nil {
			logger.Fatal().Err(err).Msg("could not unmarshal observability new relic config")
		}
		if err := k.Unmarshal("", &mainConfig.Observability.HealthChecks); err != nil {
			logger.Fatal().Err(err).Msg("could not unmarshal observability health checks config")
		}
	}

	validate := validator.New()

	err = validate.Struct(mainConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("config validation failed")
	}

	if mainConfig.Observability == nil {
		mainConfig.Observability = DefaultObservabilityConfig()
	}

	mainConfig.Observability.ServiceName = "tasker"
	mainConfig.Observability.Environment = mainConfig.Primary.Env

	if err := mainConfig.Observability.Validate(); err != nil {
		logger.Fatal().Err(err).Msg("invalid observability config")
	}

	if mainConfig.Cron == nil {
		mainConfig.Cron = DefaultCronConfig()
	}

	return mainConfig, nil
}
