package config

import (
	"fmt"
	"os"
)

var env = []string{
	"release",
	"test",
}

func GetPort() string {
	res := os.Getenv("GO_PORT")

	if res == "" {
		panic("The service port is not set")
	}

	return res
}

func GetMode() string {
	pass := false
	mode := os.Getenv("GO_MODE")

	for _, v := range env {
		if v == mode {
			pass = true
			break
		}
	}

	if pass == false {
		mode = "test"
	}

	return mode
}

func IsProduction() bool {
	result := false

	mode := GetMode()

	if mode == "release" {
		result = true
	}

	return result
}

func getKey(k string) string {
	mode := GetMode()
	result := fmt.Sprintf("%s-%s", k, mode)

	return result
}
