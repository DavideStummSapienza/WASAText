package database

import (
	"database/sql"
	"fmt"
	"time"
)

// DeleteMessage deletes a specific message from the database.
//
// This function identifies a message based on the given partner's username and timestamp.
// It verifies whether the current user has the permission to delete the message (must be the sender).
// If the conditions are met, the message is deleted from the `messages` table, and any dependent
// records are automatically deleted due to the use of `ON DELETE CASCADE` in the database schema.
//
// Parameters:
// - partnerUsername: The username of the conversation partner.
// - currentUser: The username of the current user (the one trying to delete the message).
// - messageTimestamp: The timestamp of the message to be deleted.
//
// Returns:
// - An error if the deletion fails or if the user does not have the necessary permissions.
//
// Note:
// - This function uses a SQL transaction to ensure atomicity, rolling back changes if any step fails.
func (db *appdbimpl) DeleteMessage(partnerUsername string, currentUser string, messageTimestamp time.Time) error {
	tx, err := db.c.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	var messageID int
	var sender string

	// Step 1: Find the message based on partner username and timestamp
	err = tx.QueryRow(`
        SELECT m.id, c.from_user
        FROM messages m
        INNER JOIN conversations c ON m.id = c.message_id
        WHERE c.from_user IN (?, ?)
          AND c.to_user IN (?, ?)
          AND m.created_at = ?`,
		partnerUsername, currentUser, partnerUsername, currentUser, messageTimestamp).Scan(&messageID, &sender)

	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return fmt.Errorf("message not found or no permissions to delete")
		}
		return fmt.Errorf("failed to find message: %w", err)
	}

	// Step 2: Verify permissions
	if sender != currentUser {
		tx.Rollback()
		return fmt.Errorf("user is not allowed to delete this message")
	}

	// Step 3: Delete the message (dependencies are automatically handled by CASCADE)
	_, err = tx.Exec(`DELETE FROM messages WHERE id = ?`, messageID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete message: %w", err)
	}

	// Step 4: Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
