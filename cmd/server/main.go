package main

import (
	"log"

	"go-api/internal/config"
	"go-api/internal/db"
	"go-api/internal/router"
)

func main() {
	cfg := config.Load()

	database := db.NewMySQL(cfg.DatabaseDSN)

	r := router.Setup(database)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
