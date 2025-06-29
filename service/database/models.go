package database

import "database/sql"

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Photo []byte `json:"photo,omitempty"`
}

type Group struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Photo []byte `json:"photo,omitempty"`
}

type Chat struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	CreatedAt   string         `json:"createdAt"`
	Members     []string       `json:"members"`
	LastMessage *Message       `json:"lastMessage,omitempty"`
	Messages    []Message      `json:"messages,omitempty"`
	ChatPhoto   sql.NullString `json:"chatPhoto,omitempty"`
}

type Message struct {
	Id                string   `json:"id"`
	ChatId            string   `json:"chatId"`
	SenderId          string   `json:"senderId"`
	SenderName        string   `json:"senderName"`
	Content           string   `json:"content"`
	Timestamp         string   `json:"timestamp"`
	Attachment        []byte   `json:"attachment"`
	SenderPhoto       string   `json:"senderPhoto,omitempty"`
	ReactionCount     int      `json:"reactionCount"`
	ReactingUserNames []string `json:"reactingUserNames"`
	Status            string   `json:"status"`
	ReplyTo           string   `json:"replyTo,omitempty"`
	ReplyContent      string   `json:"replyContent,omitempty"`
	ReplySenderName   string   `json:"replySenderName,omitempty"`
	ReplyAttachment   []byte   `json:"replyAttachment,omitempty"`
}

type Comment struct {
	Id       string `json:"id"`
	AuthorId string `json:"authorId"`
}

type ReadReceipt struct {
	MessageId   string  `json:"messageId"`
	UserId      string  `json:"userId"`
	DeliveredAt string  `json:"deliveredAt"`
	ReadAt      *string `json:"readAt,omitempty"`
}
