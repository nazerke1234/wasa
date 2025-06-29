package database

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

func (db *appdbimpl) CreateGroupChat(chatID string, memberIDs []string, name string, photo []byte) error {
	_, err := db.c.Exec(`
        INSERT INTO chats (id, name, type, created_at, chatPhoto)
        VALUES (?, ?, 'group', ?, ?)
    `, chatID, name, time.Now().Format(time.RFC3339), photo)
	if err != nil {
		return fmt.Errorf("error creating new chat: %w", err)
	}
	for _, memberID := range memberIDs {
		_, err = db.c.Exec(`
            INSERT INTO chat_members (chatId, userId)
            VALUES (?, ?)
        `, chatID, memberID)
		if err != nil {
			return fmt.Errorf("error adding member %s to chat_members: %w", memberID, err)
		}
	}
	return nil
}

func (db *appdbimpl) GetMyGroups(userID string) ([]Chat, error) {
	query := `
    SELECT 
        c.id,
        c.name,
        c.chatPhoto as photo
    FROM chats c
    JOIN chat_members cm ON c.id = cm.chatId
    WHERE cm.userId = ? AND c.type = 'group'
    ORDER BY c.created_at DESC;
    `
	rows, err := db.c.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching groups: %w", err)
	}
	defer rows.Close()
	var groups []Chat
	for rows.Next() {
		var group Chat
		var photo sql.NullString
		err := rows.Scan(
			&group.Id,
			&group.Name,
			&photo,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning group: %w", err)
		}
		if photo.Valid {
			group.ChatPhoto.String = base64.StdEncoding.EncodeToString([]byte(photo.String))
			group.ChatPhoto.Valid = true
		} else {
			group.ChatPhoto = sql.NullString{String: "", Valid: false}
		}
		groups = append(groups, group)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning groups: %w", err)
	}
	return groups, nil
}

func (db *appdbimpl) GetGroupInfo(groupID string) (Chat, error) {
	var group Chat
	var photo []byte
	var membersCSV sql.NullString
	err := db.c.QueryRow(`
        SELECT 
            c.id,
            c.name,
            c.chatPhoto,
            (SELECT GROUP_CONCAT(userId) FROM chat_members WHERE chatId = c.id) AS members
        FROM chats c
        WHERE c.id = ? AND c.type = 'group'`,
		groupID,
	).Scan(
		&group.Id,
		&group.Name,
		&photo,
		&membersCSV,
	)
	if err == sql.ErrNoRows {
		return Chat{}, ErrGroupDoesNotExist
	}
	if err != nil {
		return Chat{}, fmt.Errorf("error fetching group by ID: %w", err)
	}
	if len(photo) > 0 {
		group.ChatPhoto = sql.NullString{
			String: base64.StdEncoding.EncodeToString(photo),
			Valid:  true,
		}
	} else {
		group.ChatPhoto = sql.NullString{Valid: false}
	}
	if membersCSV.Valid {
		group.Members = strings.Split(membersCSV.String, ",")
	} else {
		group.Members = []string{}
	}
	return group, nil
}

func (db *appdbimpl) UpdateGroupName(groupId, newName string) error {
	res, err := db.c.Exec(`UPDATE chats SET name=? WHERE id=?`, newName, groupId)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return ErrUserDoesNotExist
	}
	return nil
}

func (db *appdbimpl) UpdateGroupPhoto(groupID string, photo []byte) error {
	var exists bool
	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM chats WHERE id=?)`, groupID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrGroupDoesNotExist
	}
	_, err = db.c.Exec(`UPDATE chats SET chatPhoto=? WHERE id=?`, photo, groupID)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) LeaveGroup(groupID, userID string) error {
	_, err := db.c.Exec(`
	DELETE FROM chat_members WHERE chatId = ? AND userId = ?
	`, groupID, userID)
	if err != nil {
		return fmt.Errorf("error leaving group: %w", err)
	}
	return nil
}

func (db *appdbimpl) AddUserToGroup(chatID string, userID string) error {
	_, err := db.c.Exec(
		"INSERT INTO chat_members (chatId, userId) VALUES (?, ?)",
		chatID, userID,
	)
	if err != nil {
		return fmt.Errorf("error adding user to group: %w", err)
	}
	return nil
}
