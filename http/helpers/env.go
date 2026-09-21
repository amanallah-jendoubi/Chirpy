package helpers

import "os"

func GetEnv(key, fallBack string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallBack
	}
	return value
}
