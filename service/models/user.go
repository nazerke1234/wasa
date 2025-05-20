package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null" json:"username"`
}

type Message struct {
	ID            uint `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	SenderID      uint           `gorm:"not null"` // Sender ID
	RecipientID   uint           `gorm:"not null"` // Recipient ID (if private message)
	ChatID        uint           `gorm:"not null"` // Chat ID (if group chat)
	Content       string         `gorm:"type:text;not null"`
	IsForwarded   bool           `gorm:"default:false"`
	ForwardedFrom *uint          `gorm:""`          // ID of the original sender (if forwarded)
	Comment       *string        `gorm:"type:text"` // Comment (reaction)
}
