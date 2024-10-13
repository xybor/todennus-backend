package config

import (
	"github.com/xybor/todennus-backend/pkg/logging"
	gormlogger "gorm.io/gorm/logger"
)

type Variable struct {
	Server         ServerVariable         `envconfig:"server"`
	Postgres       PostgresVariable       `envconfig:"postgres"`
	Redis          RedisVariable          `envconfig:"redis"`
	Authentication AuthenticationVariable `envconfig:"authentication"`
	OAuth2         OAuth2Variable         `envconfig:"oauth2"`
}

func DefaultVariable() Variable {
	return Variable{
		Server:         DefaultServerVariable(),
		Postgres:       DefaultPostgresVariable(),
		Redis:          DefaultRedisVariable(),
		Authentication: DefaultAuthenticationVariable(),
		OAuth2:         DefaultOAuth2Variable(),
	}
}

type ServerVariable struct {
	Host     string `envconfig:"host"`
	Port     int    `envconfig:"port"`
	NodeID   int    `envconfig:"nodeid"`
	LogLevel int    `envconfig:"loglevel"`
}

func DefaultServerVariable() ServerVariable {
	return ServerVariable{
		Host:     "",
		Port:     7063, // == tode
		NodeID:   1,
		LogLevel: int(logging.LevelDebug),
	}
}

type PostgresVariable struct {
	LogLevel      int    `envconfig:"loglevel"`
	Host          string `envconfig:"host"`
	Port          int    `envconfig:"port"`
	SSLMode       string `envconfig:"sslmode"`
	Timezone      string `envconfig:"timezone"`
	RetryAttempts int    `envconfig:"retry_attempts"`
	RetryInterval int    `envconfig:"retry_interval"` // in second
}

func DefaultPostgresVariable() PostgresVariable {
	return PostgresVariable{
		LogLevel:      int(gormlogger.Warn),
		Host:          "localhost",
		Port:          5432,
		SSLMode:       "disable",
		RetryAttempts: 3,
		RetryInterval: 1,
	}
}

type RedisVariable struct {
	Addr string `envconfig:"addr"`
	DB   int    `envconfig:"db"`
}

func DefaultRedisVariable() RedisVariable {
	return RedisVariable{
		Addr: "localhost:6379",
		DB:   0,
	}
}

type AuthenticationVariable struct {
	AccessTokenExpiration  int    `envconfig:"access_token_expiration"`  // in second
	RefreshTokenExpiration int    `envconfig:"refresh_token_expiration"` // in second
	IDTokenExpiration      int    `envconfig:"id_token_expiration"`      // in second
	TokenIssuer            string `envconfig:"token_issuer"`
}

func DefaultAuthenticationVariable() AuthenticationVariable {
	return AuthenticationVariable{
		AccessTokenExpiration:  60,           // 60s
		RefreshTokenExpiration: 60 * 60,      // 1h
		IDTokenExpiration:      24 * 60 * 60, // 1d
	}
}

type OAuth2Variable struct {
	ClientSecretLength int `envconfig:"client_secret_length"`
}

func DefaultOAuth2Variable() OAuth2Variable {
	return OAuth2Variable{
		ClientSecretLength: 64,
	}
}
