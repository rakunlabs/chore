package config

import "os"

func GetStoreName() string {
	envAppName := os.Getenv("APP_NAME")

	if envAppName == "" {
		Application.StoreName = Application.AppName
	} else {
		Application.StoreName = envAppName
	}

	return Application.StoreName
}
