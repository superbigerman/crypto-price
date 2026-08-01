package main

import (
	"final/internal/app"
	"log"
)

// @title Crypto Price API
// @version 1.0
// @description Сервис для получения курсов криптовалют
// @host localhost:8080
// @BasePath /

func main() {
	if err := app.Run(); err != nil {
		log.Fatal("app: %V", err)
	}
}
