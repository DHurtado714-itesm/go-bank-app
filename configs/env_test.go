package config

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

const (
	PostgresDbURI = "POSTGRES_DB_URI"
)

func TestGetString(t *testing.T) {

	getTestName := func(key string) string {
		return fmt.Sprintf("It should read the %s environment variable.", key)
	}

	tests := []struct {
		name     string
		envKey   string
		envValue string
	}{
		{
			name:     getTestName(PostgresDbURI),
			envKey:   PostgresDbURI,
			envValue: "postgresql://postgres:password@127.0.0.1:4337/some-db",
		},
	}

	for idx, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv(test.envKey, test.envValue)

			val, err := GetString(test.envKey)

			if err != nil {
				t.Fatalf("Case [%d]: Expected %s. Received error %s.", idx, test.envValue, err)
			}

			if val != test.envValue {
				t.Fatalf("Case [%d]: Expected %s. Received %s.", idx, test.envValue, val)
			}

			os.Unsetenv(test.envKey)
		})
	}
}

func TestErrGetString(t *testing.T) {

	getTestName := func(key string) string {
		return fmt.Sprintf("It should not read the %s environment variable.", key)
	}

	tests := []struct {
		name   string
		envKey string
	}{
		{
			name:   getTestName(PostgresDbURI),
			envKey: PostgresDbURI,
		},
	}

	for idx, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := GetString(test.envKey)

			if err == nil {
				t.Fatalf("Case [%d]: Expected an error but received none.", idx)
			}

			if !errors.Is(err, ErrEnvVarNotSet{key: test.envKey}) {
				t.Fatalf("Case [%d]: Expected an error of kind ErrEnvVarNotSet. Received %s.", idx, err)
			}
		})
	}
}
