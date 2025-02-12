package database

import (
	"database/sql"
	"fmt"
)

// DeleteMessage removes a message from the database if the user has permission to do so.
// It ensures that only the sender can delete their own messages.
//
// Parameters:
// - currentUser: The user attempting to delete the message.
// - messageID: The unique identifier of the message to be deleted.
//
// Returns:
// - An error if the deletion fails or if the user does not have permission.
func (db *appdbimpl) DeleteMessage(currentUser string, messageID int) error {
	tx, err := db.c.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	var sender sql.NullString

	// Step 1: Retrieve the message and check if the user has permission to delete it
	err = tx.QueryRow(`
        SELECT m.id, c.from_user
        FROM messages m
        INNER JOIN conversations c ON m.conversation_id = c.id
        WHERE m.id = ?`,
		messageID).Scan(&messageID, &sender)

	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return fmt.Errorf("message not found or no permissions to delete")
		}
		return fmt.Errorf("failed to find message: %w", err)
	}

	// Step 2: Verify that the current user is the sender of the message
	if sender.String != currentUser {
		tx.Rollback()
		return fmt.Errorf("user is not allowed to delete this message")
	}

	// Step 3: Delete the message (cascading delete will handle related records)
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
