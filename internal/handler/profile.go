package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type ProfileResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func Profile(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, _ := c.Get("user_id")
		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var res ProfileResponse
		err := db.QueryRow("SELECT first_name, last_name, email FROM users WHERE id = ?",
			userID,
		).Scan(&res.FirstName, &res.LastName, &res.Email)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch profile"})
			return
		}

		c.JSON(http.StatusOK, res)
	}
}
