package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/ini.v1"
)

type Config struct {
	Secret
	Variable
}

func Load(envPaths []string, iniPaths []string) (Config, error) {
	c := Config{}
	var err error

	c.Variable, err = loadVariables(iniPaths...)
	if err != nil {
		return c, err
	}

	c.Secret, err = loadSecret(envPaths...)
	if err != nil {
		return c, err
	}

	return c, nil
}

func loadVariables(paths ...string) (Variable, error) {
	v := DefaultVariable()

	if len(paths) > 0 {
		otherSources := []any{}
		for _, p := range paths[1:] {
			otherSources = append(otherSources, p)
		}

		initFile, err := ini.Load(paths[0], otherSources...)
		if err != nil {
			return Variable{}, err
		}

		err = initFile.MapTo(&v)
		if err != nil {
			return Variable{}, err
		}
	}

	return v, nil
}

func loadSecret(paths ...string) (Secret, error) {
	s := Secret{}

	if err := godotenv.Load(paths...); err != nil {
		return s, err
	}

	if err := envconfig.Process("postgres", &s.Postgres); err != nil {
		return s, err
	}

	if err := envconfig.Process("auth", &s.Authentication); err != nil {
		return s, err
	}

	if err := envconfig.Process("admin", &s.Admin); err != nil {
		return s, err
	}

	return s, nil
}
