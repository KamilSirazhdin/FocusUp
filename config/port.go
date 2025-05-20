package config

import "os"

func GetPort() string {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "4321"
	}
	return port
}
