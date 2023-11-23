package config

import (
	"fmt"
	"os"
)

func getEnvOptional(name, defaultValue string) string {
	value := os.Getenv(name)
	if value != "" {
		fmt.Println(value)
		return value
	}
	return defaultValue
}
