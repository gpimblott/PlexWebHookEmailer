package environment

import (
	"log"
	"os"
)

func GetEnvOrStop(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Envirnment variable %s missing", key)
	}
	return value
}

/*
Return the specified env var or the fallback if not defined.
*/
func GetEnvWithFallback(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
