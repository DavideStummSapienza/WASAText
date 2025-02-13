package database

import (
	"database/sql"
	"fmt"
)

// SendMessage sends a new message. It allows creating new 1:1 conversations
// but only allows sending messages to existing groups.
func (db *appdbimpl) SendMessage(msg NewMessage) (int, error) {
	tx, err := db.c.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}

	var conversationID int

	// 1. Check if the recipient is a group
	var isGroup bool
	err = tx.QueryRow(`SELECT COUNT(*) > 0 FROM groups WHERE groupname = ?`, msg.ToUser).Scan(&isGroup)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to check if recipient is a group: %w", err)
	}

	if isGroup {
		// 2. Ensure the sender is a member of the group
		var isMember bool
		err = tx.QueryRow(`SELECT COUNT(*) > 0 FROM group_members WHERE groupname = ? AND membername = ?`, msg.ToUser, msg.FromUser).Scan(&isMember)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to check group membership: %w", err)
		}
		if !isMember {
			tx.Rollback()
			return 0, fmt.Errorf("user %s is not a member of group %s", msg.FromUser, msg.ToUser)
		}

		// 3. Check if a conversation for the group already exists
		err = tx.QueryRow(`SELECT id FROM conversations WHERE to_group = ?`, msg.ToUser).Scan(&conversationID)
    
		if err == sql.ErrNoRows {
			// No existing conversation, create a new one
			err = tx.QueryRow(`
				INSERT INTO conversations (from_user, to_group) 
				VALUES (?, ?) RETURNING id`, msg.FromUser, msg.ToUser).Scan(&conversationID)
			if err != nil {
				tx.Rollback()
				return 0, fmt.Errorf("failed to create new group conversation: %w", err)
			}
		} else if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to fetch group conversation: %w", err)
		}

	} else {
		// 4. Fetch or create a new 1:1 conversation
		err = tx.QueryRow(`
			SELECT id FROM conversations 
			WHERE from_user = ? AND to_user = ?`,
			msg.FromUser, msg.ToUser).Scan(&conversationID)

		if err == sql.ErrNoRows {
			err = tx.QueryRow(`
				INSERT INTO conversations (from_user, to_user) VALUES (?, ?) RETURNING id`,
				msg.FromUser, msg.ToUser).Scan(&conversationID)
			if err != nil {
				tx.Rollback()
				return 0, fmt.Errorf("failed to create new conversation: %w", err)
			}
		} else if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to fetch conversation: %w", err)
		}
	}

	// 5. Insert new message 
	var messageID int
	err = tx.QueryRow(`
		INSERT INTO messages (content, is_photo, photo_url, is_forwarded, created_at, conversation_id) 
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, ?) RETURNING id`,
		msg.Content, msg.IsPhoto, msg.PhotoURL, msg.IsForwarded, conversationID).Scan(&messageID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert new message: %w", err)
	}
	

	// 6. Mark message as "unread" and "unreceived"for the recipients
	if isGroup {
		// Mark as unread and unreceived for all group members (excluding the sender)
		_, err = tx.Exec(`
			INSERT INTO message_status (message_id, user_id, received, read) 
			SELECT ?, membername, FALSE, FALSE FROM group_members WHERE groupname = ? AND membername != ?`,
			messageID, msg.ToUser, msg.FromUser)
	} else {
		// Mark as unread for a 1:1 recipient
		_, err = tx.Exec(`
			INSERT INTO message_status (message_id, user_id, received, read) 
			VALUES (?, ?, FALSE, FALSE)`, messageID, msg.ToUser)
	}
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert message status: %w", err)
	}

	// 7. Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return messageID, nil
}
