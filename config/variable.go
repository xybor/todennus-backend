package config

import (
	"github.com/xybor/todennus-backend/pkg/logging"
	gormlogger "gorm.io/gorm/logger"
)

type Variable struct {
	Server         ServerVariable         `ini:"server"`
	Postgres       PostgresVariable       `ini:"postgres"`
	Authentication AuthenticationVariable `ini:"authentication"`
	OAuth2         OAuth2Variable         `ini:"oauth2"`
}

func DefaultVariable() Variable {
	return Variable{
		Server:         DefaultSystemVariable(),
		Postgres:       DefaultPostgresVariable(),
		Authentication: DefaultAuthenticationVariable(),
		OAuth2:         DefaultOAuth2Variable(),
	}
}

type ServerVariable struct {
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	NodeID   int    `ini:"nodeid"`
	LogLevel int    `ini:"loglevel"`
}

func DefaultSystemVariable() ServerVariable {
	return ServerVariable{
		Host:     "",
		Port:     7063, // == tode
		NodeID:   1,
		LogLevel: int(logging.LevelDebug),
	}
}

type PostgresVariable struct {
	LogLevel      int    `ini:"loglevel"`
	Host          string `ini:"host"`
	Port          int    `ini:"port"`
	SSLMode       string `ini:"sslmode"`
	Timezone      string `ini:"timezone"`
	RetryAttempts int    `ini:"retry_attempts"`
	RetryInterval int    `ini:"retry_interval"` // in second
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

type AuthenticationVariable struct {
	AccessTokenExpiration  int    `ini:"access_token_expiration"`  // in second
	RefreshTokenExpiration int    `ini:"refresh_token_expiration"` // in second
	IDTokenExpiration      int    `ini:"id_token_expiration"`      // in second
	TokenIssuer            string `ini:"token_issuer"`
}

func DefaultAuthenticationVariable() AuthenticationVariable {
	return AuthenticationVariable{
		AccessTokenExpiration:  60,           // 60s
		RefreshTokenExpiration: 60 * 60,      // 1h
		IDTokenExpiration:      24 * 60 * 60, // 1d
	}
}

type OAuth2Variable struct {
	ClientSecretLength int `ini:"client_secret_length"`
}

func DefaultOAuth2Variable() OAuth2Variable {
	return OAuth2Variable{
		ClientSecretLength: 32,
	}
}
