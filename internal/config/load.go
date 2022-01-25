package config

import "os"

func GetLoadName() string {
	envAppName := os.Getenv("APP_NAME")

	if envAppName == "" {
		LoadName = AppName
	} else {
		LoadName = envAppName
	}

	return LoadName
}
