package controllers

import (
	"net/http"
	"wasa/backend/internal/utils"
	"wasa/backend/service/database"
	"wasa/backend/service/models"

	"github.com/gin-gonic/gin"
)

// Get profile information (authorized user)
func GetProfile(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user database.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"photo":    user.Photo,
	})
}

// Update username
func UpdateUsername(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user database.User
	// Find user by ID, not username
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input struct {
		NewUsername string `json:"new_username"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if len(input.NewUsername) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be at least 2 characters long"})
		return
	}

	// Check if the new username is unique
	var existingUser database.User
	if err := database.DB.Where("username = ?", input.NewUsername).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username is already taken"})
		return
	}

	// Update username
	user.Username = input.NewUsername
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving username"})
		return
	}

	// Generate new token
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	// Send new token to the user
	c.JSON(http.StatusOK, gin.H{
		"message":       "Username updated",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Update profile photo (expects a URL)
func UpdateProfilePhoto(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user database.User
	// Find user by ID, not username
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input struct {
		PhotoURL string `json:"photo_url"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Update profile photo
	user.Photo = input.PhotoURL
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating profile photo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile photo updated"})
}

// Search users by name
func SearchUsers(c *gin.Context) {
	// Get name for search
	query := c.Query("name")
	if len(query) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter at least 3 letters"})
		return
	}

	// Perform user search
	var users []models.User
	err := database.DB.Where("username LIKE ?", "%"+query+"%").
		Order("username ASC"). // Add sorting
		Find(&users).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching"})
		return
	}

	c.JSON(http.StatusOK, users)
}
