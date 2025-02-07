package controllers

import (
	"net/http"
	"strconv"
	"test-golang/config"
	"test-golang/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUsers - Get all users with pagination
func GetUsers(c *gin.Context) {
	var users []models.User

	// Get 'page' and 'results' query parameters
	page := c.DefaultQuery("page", "1")        // Default to page 1
	results := c.DefaultQuery("results", "10") // Default to 10 results per page

	// Convert query parameters to integers
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	resultsInt, err := strconv.Atoi(results)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid results number"})
		return
	}

	// Calculate the offset for pagination
	offset := (pageInt - 1) * resultsInt

	// Query with pagination
	if err := config.DB.Limit(resultsInt).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving users"})
		return
	}

	// Return the paginated result
	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	// Validasi UUID
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// Preload company agar data lengkap
	if err := config.DB.Preload("Company").First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// Create a new user
func CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate UUID untuk user baru
	user.ID = uuid.New()
	user.Company.ID = uuid.New()
	user.Company.UserID = user.ID // Set relasi One-to-One

	config.DB.Create(&user)
	c.JSON(http.StatusCreated, user)
}

// Update user by UUID
func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	// Validasi UUID
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	if err := config.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pastikan ID tidak berubah
	user.ID, _ = uuid.Parse(id)

	config.DB.Save(&user)
	c.JSON(http.StatusOK, user)
}

// Delete user by UUID
func DeleteUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	// Validasi UUID
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// Cek apakah user ada
	if err := config.DB.Preload("Company").First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	config.DB.Delete(&user) // Akan otomatis menghapus Company karena OnDelete:CASCADE
	c.JSON(http.StatusOK, gin.H{"message": "User and related Company deleted"})
}
