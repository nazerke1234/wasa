package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	database2 "wasa/backend/service/database"

	"github.com/gin-gonic/gin"
)

// Logic for adding users to the group
func AddUserToGroup(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized or invalid user_id"})
		return
	}

	groupID, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	var input struct {
		UserIDs []uint `json:"user_ids"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	var chat database2.Chat
	if err := database2.DB.Where("id = ? AND is_group = ?", groupID, true).First(&chat).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	for _, userID := range input.UserIDs {
		chatUser := database2.ChatUser{
			ChatID: uint(groupID),
			UserID: userID,
		}

		if err := database2.DB.Create(&chatUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding user to the group"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Users added to the group"})
}

// Leaving the group
func LeaveGroup(c *gin.Context) {
	// Get user_id from the context
	userID, exists := c.Get("user_id")
	if !exists || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get chatId from URL parameters
	chatID, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	fmt.Println("CHATID", chatID)

	// Check if the group exists
	var chat database2.Chat
	if err := database2.DB.Where("id = ? AND is_group = ?", chatID, true).First(&chat).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Remove user from the group
	if err := database2.DB.Where("chat_id = ? AND user_id = ?", chatID, userID).Delete(&database2.ChatUser{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error leaving the group"})
		return
	}

	// If the user leaves and it's the last user, we can delete the chat
	var chatUsers []database2.ChatUser
	if err := database2.DB.Where("chat_id = ?", chatID).Find(&chatUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking group members"})
		return
	}

	if len(chatUsers) == 0 {
		// Delete the chat if there are no more members
		if err := database2.DB.Delete(&chat).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting chat"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "You have successfully left the group"})
}

// Getting group information
func GetGroupInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	groupID, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var chat database2.Chat
	if err := database2.DB.Where("id = ? AND is_group = ?", groupID, true).First(&chat).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Get the number of members in the group
	var userCount int64
	database2.DB.Model(&database2.ChatUser{}).Where("chat_id = ?", groupID).Count(&userCount)

	// Return the group name and number of members
	c.JSON(http.StatusOK, gin.H{
		"name":        chat.Name,
		"photo":       chat.Photo,
		"users_count": userCount,
	})
}

// Updating group name
func SetGroupName(c *gin.Context) {
	// Get user_id from the context
	userID, exists := c.Get("user_id")
	if !exists || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized or invalid user_id"})
		return
	}

	// Get chatId from URL parameters
	chatID, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	// Get the new group name from the request body
	var input struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check if the group exists with the given chatId
	var chat database2.Chat
	if err := database2.DB.Where("id = ? AND is_group = ?", chatID, true).First(&chat).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Update the group name
	chat.Name = input.Name
	if err := database2.DB.Save(&chat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating group name"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group name updated", "chat_name": chat.Name})
}

// Updating group photo
func UpdateGroupPhoto(c *gin.Context) {
	// Get user_id from the context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get chatId from URL parameters
	chatID, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	// Check if the group exists
	var chat database2.Chat
	if err := database2.DB.Where("id = ? AND is_group = ?", chatID, true).First(&chat).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Check if the user is a member of the group
	var chatUser database2.ChatUser
	if err := database2.DB.Where("chat_id = ? AND user_id = ?", chatID, userID).First(&chatUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not a member of this group"})
		return
	}

	// Get the new photo URL
	var input struct {
		PhotoURL string `json:"photo_url"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Update the group photo
	chat.Photo = input.PhotoURL
	if err := database2.DB.Save(&chat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating group photo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group photo updated"})
}
