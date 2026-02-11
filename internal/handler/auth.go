package handler

import (
	"net/http"

	"go-api/internal/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Email        string `json:"email"`
}

type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

var jwtSecret = []byte("SUPER_SECRET_KEY")

func Login(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user model.User
		err := db.QueryRow("SELECT id, email, password_hash FROM users WHERE email = ? AND status = 1",
			req.Email,
		).Scan(&user.ID, &user.Email, &user.PasswordHash)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if err := bcrypt.CompareHashAndPassword(
			[]byte(user.PasswordHash),
			[]byte(req.Password),
		); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
			"exp":     time.Now().Add(15 * time.Minute).Unix(),
		})

		accessString, err := accessToken.SignedString(jwtSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
			"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		})

		refreshString, err := refreshToken.SignedString(jwtSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		_, err = db.Exec("INSERT INTO refresh_tokens (user_id, token, expired_at) VALUES (?, ?, ?)", user.ID, refreshString, time.Now().Add(7*24*time.Hour))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, LoginResponse{
			AccessToken:  accessString,
			RefreshToken: refreshString,
			Email:        user.Email,
		})
	}
}

func Register(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// check email if exists
		var exists int

		err := db.QueryRow("SELECT 1 FROM users WHERE email = ? LIMIT 1", req.Email).Scan(&exists)

		if err == nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"error": "Email already registered"})
			return
		}

		hashed, err := bcrypt.GenerateFromPassword(
			[]byte(req.Password),
			bcrypt.DefaultCost,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to hash password",
			})
			return
		}

		_, err = db.Exec(
			"INSERT INTO users (first_name, last_name, email, password_hash) VALUES (?, ?, ?, ?)",
			req.FirstName, req.LastName, req.Email, string(hashed),
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to register user",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User registered successfully",
		})
	}
}
