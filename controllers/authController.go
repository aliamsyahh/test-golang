package controllers

import (
	"net/http"
	"time"

	"test-golang/config"
	"test-golang/models"
	"test-golang/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Struct untuk menangkap request login
type LoginRequest struct {
	Email string `json:"email" binding:"required"`
}

// Login handler untuk autentikasi user
func Login(c *gin.Context) {
	var input LoginRequest

	// Bind JSON request ke struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	// Cek apakah user ada di database
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		}
		return
	}

	// Generate JWT token
	token, err := service.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token"})
		return
	}

	// Set token dalam cookie (secure HttpOnly)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session_id",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user,
		"token":   token,
	})
}

// Protected route that requires valid JWT.
func Protected(c *gin.Context) {
	// Retrieve the cookie containing the JWT.
	cookie, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "No session found"})
		return
	}

	// Validate the JWT.
	claims, err := service.ValidateJWT(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}

	// Return protected data.
	c.JSON(http.StatusOK, gin.H{
		"message": "Protected data",
		"user":    claims.Username,
	})
}
