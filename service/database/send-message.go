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

	// Defer rollback to ensure transaction cleanup in case of errors
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w, original error: %w", rollbackErr, err)
			}
		}
	}()

	var conversationID int

	// 1. Check if the recipient is a group
	var isGroup bool
	err = tx.QueryRow(`SELECT COUNT(*) > 0 FROM groups WHERE groupname = ?`, msg.ToUser).Scan(&isGroup)
	if err != nil {
		return 0, fmt.Errorf("failed to check if recipient is a group: %w", err)
	}

	if isGroup {
		// 2. Ensure the sender is a member of the group
		var isMember bool
		err = tx.QueryRow(`SELECT COUNT(*) > 0 FROM group_members WHERE groupname = ? AND membername = ?`, msg.ToUser, msg.FromUser).Scan(&isMember)
		if err != nil {
			return 0, fmt.Errorf("failed to check group membership: %w", err)
		}
		if !isMember {
			return 0, fmt.Errorf("user %s is not a member of group %s", msg.FromUser, msg.ToUser)
		}

		// 3. Check if a conversation for the group already exists
		err = tx.QueryRow(`SELECT id FROM conversations WHERE groupname = ?`, msg.ToUser).Scan(&conversationID)

		if err == sql.ErrNoRows {
			// No existing conversation, create a new one
			err = tx.QueryRow(`
				INSERT INTO conversations (groupname) 
				VALUES (?) RETURNING id`, msg.ToUser).Scan(&conversationID)
			if err != nil {
				return 0, fmt.Errorf("failed to create new group conversation: %w", err)
			}
		} else if err != nil {
			return 0, fmt.Errorf("failed to fetch group conversation: %w", err)
		}

	} else {
		// 4. Fetch or create a new 1:1 conversation
		err = tx.QueryRow(`
			SELECT id FROM conversations 
			WHERE (user1 = ? AND user2 = ?) OR (user2 = ? AND user1 = ?)`,
			msg.FromUser, msg.ToUser, msg.FromUser, msg.ToUser).Scan(&conversationID)

		if err == sql.ErrNoRows {
			err = tx.QueryRow(`
				INSERT INTO conversations (user1, user2) VALUES (?, ?) RETURNING id`,
				msg.FromUser, msg.ToUser).Scan(&conversationID)
			if err != nil {
				return 0, fmt.Errorf("failed to create new conversation: %w", err)
			}
		} else if err != nil {
			return 0, fmt.Errorf("failed to fetch conversation: %w", err)
		}
	}

	// 5. Insert new message
	var messageID int
	err = tx.QueryRow(`
		INSERT INTO messages (content, sender, is_photo, is_forwarded, created_at, conversation_id) 
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, ?) RETURNING id`,
		msg.Content, msg.FromUser, msg.IsPhoto, msg.IsForwarded, conversationID).Scan(&messageID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert new message: %w", err)
	}

	// 6. Mark message as "unread" and "unreceived" for the recipients
	if isGroup {
		_, err = tx.Exec(`
			INSERT INTO message_status (message_id, user_id, received, read) 
			SELECT ?, membername, FALSE, FALSE FROM group_members WHERE groupname = ? AND membername != ?`,
			messageID, msg.ToUser, msg.FromUser)
	} else {
		_, err = tx.Exec(`
			INSERT INTO message_status (message_id, user_id, received, read) 
			VALUES (?, ?, FALSE, FALSE)`, messageID, msg.ToUser)
	}
	if err != nil {
		return 0, fmt.Errorf("failed to insert message status: %w", err)
	}

	// 7. Commit the transaction
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return messageID, nil
}
