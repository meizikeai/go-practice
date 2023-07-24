package config

import (
	"os"
)

func getMode() string {
	mode := os.Getenv("GO_MODE")

	if mode == "" {
		mode = "test"
	}

	return mode
}

func isProduction() bool {
	result := false

	mode := os.Getenv("GO_MODE")

	if mode == "release" {
		result = true
	}

	return result
}
