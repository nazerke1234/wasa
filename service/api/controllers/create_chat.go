package controllers

import (
	"fmt"
	"net/http"
	database2 "wasa/backend/service/database"

	"github.com/gin-gonic/gin"
)

// CreatePersonalChat creates a personal chat between the user and another user.
func CreatePersonalChat(c *gin.Context) {
	// Get user_id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		UserID uint `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check if the second user exists
	var otherUser database2.User
	if err := database2.DB.Where("id = ?", input.UserID).First(&otherUser).Error; err != nil {
		fmt.Println("User not found with ID:", input.UserID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// If the user wants to create a chat with themselves
	if userID == input.UserID {
		// Check if a chat with the name "me" already exists
		var existingChat database2.Chat
		err := database2.DB.Where("name = ? AND id IN (SELECT chat_id FROM chat_users WHERE user_id = ?)", "me", userID).First(&existingChat).Error

		// If chat exists, return its ID
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"message": "Chat already exists", "chat_id": existingChat.ID})
			return
		}

		// Create a chat named "me"
		chat := database2.Chat{
			Name:    "me",
			IsGroup: false, // This is a personal chat
		}

		if err := database2.DB.Create(&chat).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating chat"})
			return
		}

		// Add user to the chat
		chatUser := database2.ChatUser{
			ChatID: chat.ID,
			UserID: userID.(uint),
		}

		if err := database2.DB.Create(&chatUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding user to chat"})
			return
		}

		// Return the created chat's ID
		c.JSON(http.StatusOK, gin.H{"message": "Chat with yourself created", "chat_id": chat.ID})
		return
	}

	// Check if a chat with the specified user exists in the chat_users table
	var existingChat database2.ChatUser
	err := database2.DB.Where("(user_id = ? AND chat_id IN (SELECT chat_id FROM chat_users WHERE user_id = ?)) AND chat_id IN (SELECT id FROM chats WHERE is_group = false)", userID, input.UserID).First(&existingChat).Error

	// If chat exists, return its ID
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Chat already exists", "chat_id": existingChat.ChatID})
		return
	}

	var currentUser database2.User
	err = database2.DB.Where("id = ?", userID).First(&currentUser).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user data"})
		return
	}

	// Create a new chat
	chat := database2.Chat{
		Name:    currentUser.Username + " and " + otherUser.Username,
		IsGroup: false, // This is a personal chat
	}

	if err := database2.DB.Create(&chat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating chat"})
		return
	}

	// Add users to the chat
	chatUser1 := database2.ChatUser{
		ChatID: chat.ID,
		UserID: userID.(uint),
	}

	chatUser2 := database2.ChatUser{
		ChatID: chat.ID,
		UserID: input.UserID,
	}

	if err := database2.DB.Create(&chatUser1).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding user to chat"})
		return
	}

	if err := database2.DB.Create(&chatUser2).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding user to chat"})
		return
	}

	// Return the created chat's ID and the other user's name
	c.JSON(http.StatusOK, gin.H{"message": "Chat created", "chat_id": chat.ID, "user_name": otherUser.Username})
}
