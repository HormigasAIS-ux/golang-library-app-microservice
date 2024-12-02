package main

import (
	"auth_service/config"
)

func init() {
	config.InitEnv("./.env")
	config.ConfigureLogger()
}

func main() {
}
