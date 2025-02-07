package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"test-golang/config"
	"test-golang/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

// GoogleLogin will initiate the OAuth flow.
func GoogleLogin(c *gin.Context) {
	// Debug log untuk memastikan client_id sudah terisi dengan benar
	log.Println("Google OAuth Config:", config.GoogleOauthConfig)

	url := config.GoogleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

// GoogleCallback is called after the user is redirected back from Google.
func GoogleCallback(c *gin.Context) {
	// Retrieve the authorization code from the query parameters
	code := c.DefaultQuery("code", "")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	// Exchange the authorization code for an access token
	token, err := config.GoogleOauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Use the token to get user info from Google
	client := config.GoogleOauthConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()
	// Parse the user information
	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	// Check if the user exists in the database
	var user models.User
	if err := config.DB.Where("email = ?", userInfo.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// User not found, create a new user
			user = models.User{
				ID:    uuid.New(),
				Name:  userInfo.Name,
				Email: userInfo.Email,
			}
			config.DB.Create(&user)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	}

	// Generate JWT token for the user
	jwtToken, err := config.GenerateJWT(user.ID.String(), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create JWT token"})
		return
	}

	// Set the JWT token in the cookie (valid for 24 hours)
	c.SetCookie("jwt", jwtToken, 3600*24, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// Logout will delete the JWT cookie to log the user out.
func Logout(c *gin.Context) {
	// Clear the JWT cookie
	c.SetCookie("jwt", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
