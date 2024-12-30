package database

import (
	"database/sql"
	"fmt"
)

// SendMessage sends a new message. If the conversation doesn't exist, it creates it first.
func (db *appdbimpl) SendMessage(fromUser, toUser, messageContent string, isPhoto bool, photoURL string) (int, error) {
	tx, err := db.c.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}

	var conversationID int

	// Step 1: Check if the conversation exists
	err = tx.QueryRow(`
        SELECT id
        FROM conversations
        WHERE (from_user = ? AND to_user = ?)
           OR (from_user = ? AND to_user = ?)`,
		fromUser, toUser, toUser, fromUser).Scan(&conversationID)

	if err != nil {
		if err == sql.ErrNoRows {
			// Step 2: Create a new conversation if it doesn't exist
			err = tx.QueryRow(`
                INSERT INTO conversations (from_user, to_user, message_id)
                VALUES (?, ?, NULL)
                RETURNING id`,
				fromUser, toUser).Scan(&conversationID)
			if err != nil {
				tx.Rollback()
				return 0, fmt.Errorf("failed to create new conversation: %w", err)
			}
		} else {
			tx.Rollback()
			return 0, fmt.Errorf("failed to check conversation existence: %w", err)
		}
	}

	// Step 3: Insert the new message into the messages table
	var messageID int
	err = tx.QueryRow(`
        INSERT INTO messages (content, is_photo, photo_url, created_at)
        VALUES (?, ?, ?, CURRENT_TIMESTAMP)
        RETURNING id`,
		messageContent, isPhoto, photoURL).Scan(&messageID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert new message: %w", err)
	}

	// Step 4: Link the new message to the conversation
	_, err = tx.Exec(`
        UPDATE conversations
        SET message_id = ?
        WHERE id = ?`,
		messageID, conversationID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to link message to conversation: %w", err)
	}

	// Step 5: Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return messageID, nil
}
