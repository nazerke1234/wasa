package controllers

import (
	"net/http"
	database2 "wasa/service/database"

	"github.com/gin-gonic/gin"
)

// CreateGroupChat creates a new group chat
func CreateGroupChat(c *gin.Context) {
	// Get user_id from context
	userID, exists := c.Get("user_id")
	if !exists || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized or invalid user_id"})
		return
	}

	// Convert userID from interface{} to uint
	userIDInt, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id"})
		return
	}

	var input struct {
		UserIDs []uint `json:"user_ids"` // List of users for the group chat
	}

	// Bind incoming data
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check that the user list is not empty
	if len(input.UserIDs) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least 2 users are required for a group chat"})
		return
	}

	// Add the current user to the list if not already present
	if !contains(input.UserIDs, userIDInt) {
		input.UserIDs = append(input.UserIDs, userIDInt)
	}

	// Create a new group chat
	chat := database2.Chat{
		Name:    "Group Chat", // The group name can be customized
		IsGroup: true,
	}

	// Create the chat in the database
	if err := database2.DB.Create(&chat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating the chat"})
		return
	}

	// Add users to the chat
	for _, userID := range input.UserIDs {
		if userID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var user database2.User
		// Find the user by ID
		if err := database2.DB.Where("id = ?", userID).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Add the relationship between the user and the chat
		chatUser := database2.ChatUser{
			ChatID: chat.ID,
			UserID: userID,
		}

		// Save the relationship in the chat_users table
		if err := database2.DB.Create(&chatUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding user to the chat"})
			return
		}
	}

	// Return the ID of the created chat
	c.JSON(http.StatusOK, gin.H{"message": "Group chat created", "chat_id": chat.ID})
}

// Check if an element exists in a slice
func contains(slice []uint, item uint) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
