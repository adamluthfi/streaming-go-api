package main

import (
	"go-api/internal/db"
	"go-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	dsn := "app_user:123qweasd@tcp(localhost:3306)/stream_app_db?parseTime=true&charset=utf8mb4"
	db := db.NewMySQL(dsn)
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	r.GET("/ping", handler.Ping)

	r.POST("/login", handler.Login(db))

	r.POST("/register", handler.Register(db))

	r.Run(":8080")
}
