package controllers

import (
	"net/http"
	database2 "wasa/service/database"

	"github.com/gin-gonic/gin"
)

// Get the list of the user's conversations via the relationship with chat_users
func GetMyConversations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var chats []database2.Chat
	// Query through the relationship with chat_users table
	result := database2.DB.Joins("JOIN chat_users ON chat_users.chat_id = chats.id").
		Where("chat_users.user_id = ?", userID).
		Find(&chats)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving chats"})
		return
	}

	var currentUser database2.User
	err := database2.DB.First(&currentUser, "id = ?", userID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user data"})
		return
	}

	var chatData []map[string]interface{}
	for _, chat := range chats {
		chatInfo := map[string]interface{}{
			"ID":      chat.ID,
			"Name":    chat.Name,
			"IsGroup": chat.IsGroup,
		}

		if chat.IsGroup {
			chatInfo["PhotoUrl"] = chat.Photo
		} else {
			if chat.Name == "me" {
				chatInfo["PhotoUrl"] = currentUser.Photo // Adding current user's photo
			} else {
				// For personal chat, find the other participant (excluding current user)
				var otherUser database2.User
				err := database2.DB.Joins("JOIN chat_users cu ON cu.chat_id = ? AND cu.user_id != ?", chat.ID, userID).
					Joins("JOIN users u ON cu.user_id = u.id").
					Select("u.photo, u.username"). // Ensuring we select the right fields
					First(&otherUser).Error

				if err != nil {
					// If user not found, set default photo
					chatInfo["PhotoUrl"] = "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQtr3FkBqr_ikozzp2wuBqXRrYlZN_q4SNqcQ&s"
				} else {
					chatInfo["PhotoUrl"] = otherUser.Photo
					chatInfo["OtherUserName"] = otherUser.Username // Adding the other user's name
				}
			}
		}

		chatData = append(chatData, chatInfo)
	}

	c.JSON(http.StatusOK, chatData)
}

// Get messages from a specific chat
func GetConversation(c *gin.Context) {
	chatID := c.Param("chat_id")

	var messages []database2.Message
	database2.DB.Where("chat_id = ?", chatID).Find(&messages)

	c.JSON(http.StatusOK, messages)
}

// Send a message to the chat (for both personal and group chats)
func SendMessage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		ChatID  uint   `json:"chat_id"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Get the chat by ID
	var chat database2.Chat
	if err := database2.DB.First(&chat, input.ChatID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	// Create a message and save it
	message := database2.Message{
		SenderID: userID.(uint),
		ChatID:   input.ChatID,
		Content:  input.Content,
		// RecipientID: 0,
	}

	if err := database2.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending message"})
		return
	}

	// If it's a personal chat, the message will be sent only to the recipient
	// If it's a group chat, it will be sent to all participants

	// If the chat is not a group, find the recipient
	if !chat.IsGroup {
		var chatUser database2.ChatUser
		err := database2.DB.Where("chat_id = ? AND user_id != ?", input.ChatID, userID).First(&chatUser).Error

		if err != nil {
			// Check if this is a chat with the user themselves
			var count int64
			database2.DB.Model(&database2.ChatUser{}).Where("chat_id = ?", input.ChatID).Count(&count)

			if count == 1 { // Only one participant in the chat, meaning it's a chat with the user themselves
				c.JSON(http.StatusCreated, message)
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Recipient not found"})
			return
		}

		// message.RecipientID = chatUser.UserID
		database2.DB.Save(&message)
	}

	c.JSON(http.StatusCreated, message)
}

// Forward a message
func ForwardMessage(c *gin.Context) {
	// Get the current user's ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Read data from the JSON request
	var input struct {
		MessageID uint   `json:"message_id"`
		ChatName  string `json:"chat_name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Find the original message
	var originalMessage database2.Message
	if err := database2.DB.First(&originalMessage, input.MessageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Find the chat by name
	var chat database2.Chat
	if err := database2.DB.Where("name = ?", input.ChatName).First(&chat).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	// Check if the user is in this chat
	var chatUser database2.ChatUser
	if err := database2.DB.Where("chat_id = ? AND user_id = ?", chat.ID, userID).First(&chatUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not a member of this chat"})
		return
	}

	// Create a new message (forwarded)
	newMessage := database2.Message{
		SenderID:      userID.(uint),
		ChatID:        chat.ID, // New chat
		Content:       originalMessage.Content,
		IsForwarded:   true,
		ForwardedFrom: &originalMessage.SenderID,
	}

	// Save to DB
	if err := database2.DB.Create(&newMessage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error forwarding"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "✅ Message forwarded!"})
}

// Add a reaction (comment)
func CommentMessage(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	messageID := c.Param("message_id")
	var input struct {
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	var message database2.Message
	if err := database2.DB.First(&message, messageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	message.Comment = &input.Comment
	if err := database2.DB.Save(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment added"})
}

// Remove a comment
func UncommentMessage(c *gin.Context) {
	messageID := c.Param("message_id")

	var message database2.Message
	if err := database2.DB.First(&message, messageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Check if there is a comment
	if message.Comment == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No comment to remove"})
		return
	}

	message.Comment = nil
	if err := database2.DB.Save(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment removed"})
}

// Delete a message
func DeleteMessage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	messageID := c.Param("message_id")

	var message database2.Message
	if err := database2.DB.First(&message, messageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Check if the current user is the sender of the message
	if message.SenderID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You cannot delete this message"})
		return
	}

	database2.DB.Delete(&message)
	c.JSON(http.StatusOK, gin.H{"message": "Message deleted"})
}
