package config

import (
	"fmt"
	"os"
)

type ErrEnvVarNotSet struct {
	key string
}

func (err ErrEnvVarNotSet) Error() string {
	return fmt.Sprintf("Environment variable with key %s was not set", err.key)
}

// GetString reads a string value from the environment variables. Returns an error if the variable with the provided key is not set.
func GetString(key string) (string, error) {
	value, ok := os.LookupEnv(key)

	if !ok {
		return value, ErrEnvVarNotSet{key: key}
	}

	return value, nil
}

// GetString reads a string value from the environment variables. If the enviroment variable is not set the provided default value will be used.
func GetStringOrDefault(key string, defaultValue string) string {
	value, err := GetString(key)

	if err != nil {
		value = defaultValue
	}

	return value
}