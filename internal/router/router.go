package router

import (
	"go-api/internal/handler"
	"go-api/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Setup(db *sqlx.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/health", handler.Health)
	r.GET("/ping", handler.Ping)
	r.POST("/register", handler.Register(db))
	r.POST("/login", handler.Login(db))

	auth := r.Group("/")

	auth.Use(middleware.JWTAuth())
	{
		auth.GET("/userID", handler.User())
		auth.GET("/profile", handler.Profile(db))
	}

	return r
}
