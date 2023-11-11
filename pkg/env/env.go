package env

import (
	"os"
)

func Get(envName string) string {
	return os.Getenv(envName)
}

func GetOrDefault(envName string, defaultValue string) string {
	val, found := os.LookupEnv(envName)
	if !found {
		return defaultValue
	}
	return val
}

func GetOrPanic(envName string) string {
	val, found := os.LookupEnv(envName)
	if !found {
		panic("missing env var: " + envName)
	}
	return val
}

func SetEnvIfNotDefined(key, value string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, value)
	}
}
