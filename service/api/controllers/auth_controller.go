package controllers

import (
	"errors"
	"net/http"
	"wasa/internal/utils"
	database2 "wasa/service/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Request structures
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	var user database2.User
	result := database2.DB.Where("username = ?", input.Username).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// User not found — create a new one
		user = database2.User{Username: input.Username}
		database2.DB.Create(&user)
	}

	// Create JWT token
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":       user.ID,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Handler for refreshing the access token
func RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Validate the refresh token
	claims, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Extract username and user_id from the token
	userID := claims.UserID
	username := claims.Username

	// Generate a new access token
	accessToken, _, err := utils.GenerateTokens(userID, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}
