package database

import (
	"fmt"
	"gorm.io/gorm"
)

// User model
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null"`
	Photo    string // Profile photo URL
}

// Chat model (group or private)
type Chat struct {
	gorm.Model
	Name    string
	IsGroup bool   // Flag indicating if the chat is a group
	Photo   string // Group photo
	Members []User `gorm:"many2many:chat_users;"` // Many-to-many relationship
}

// Message model
type Message struct {
	gorm.Model
	ChatID        uint
	SenderID      uint
	Content       string
	RecipientID   uint
	Comment       *string
	IsForwarded   bool
	ForwardedFrom *uint
}

type ChatUser struct {
	ChatID uint `gorm:"primaryKey"`
	UserID uint `gorm:"primaryKey"`
}

// Migration function
func MigrateDB() {
	// Perform migration to create tables
	err := DB.AutoMigrate(&User{}, &Chat{}, &Message{}, &ChatUser{})
	if err != nil {
		fmt.Println("Error during migration:", err)
	} else {
		fmt.Println("Migration completed successfully")
	}
}
