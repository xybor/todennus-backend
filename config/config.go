package config

import (
	"reflect"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Secret
	Variable
}

func Load(paths ...string) (Config, error) {
	if len(paths) > 0 {
		if err := godotenv.Load(paths...); err != nil {
			return Config{}, err
		}
	}

	c := Config{}
	c.Variable = DefaultVariable()

	if err := load(&c.Variable); err != nil {
		return Config{}, err
	}

	if err := load(&c.Secret); err != nil {
		return Config{}, err
	}

	return c, nil
}

func load[T any](obj *T) error {
	sType := reflect.TypeOf(obj).Elem()
	sValue := reflect.ValueOf(obj).Elem()
	for i := range sType.NumField() {
		field := sType.Field(i)
		prefix := field.Tag.Get("envconfig")
		if prefix == "" {
			prefix = strings.ToLower(field.Name)
		}

		if err := envconfig.Process(prefix, sValue.FieldByName(field.Name).Addr().Interface()); err != nil {
			return err
		}
	}

	return nil
}
