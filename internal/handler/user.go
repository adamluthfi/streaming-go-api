package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func User() gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, _ := c.Get("user_id")

		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
			"message": "User is authenticated",
		})
	}
}
