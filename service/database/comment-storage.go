package database

import (
	"fmt"
)

func (db *appdbimpl) CommentMessage(commentID, messageID, authorID string) error {
	query := `INSERT INTO comments (id, messageId, authorId) VALUES (?, ?, ?)`
	_, err := db.c.Exec(query, commentID, messageID, authorID)
	if err != nil {
		return fmt.Errorf("failed to insert comment: %w", err)
	}
	return nil
}

func (db *appdbimpl) UncommentMessage(messageID, authorID string) error {
	query := `DELETE FROM comments WHERE messageId = ? AND authorId = ?`
	_, err := db.c.Exec(query, messageID, authorID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}
