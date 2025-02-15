package database

import "fmt"

// ChangeUsername updates the username of an existing user in the database.
// This function assumes that the username validation has already been performed at a higher level (API layer).
func (db *appdbimpl) ChangeUsername(oldUsername, newUsername string) error {
	// Check if the new username is the same as the old one
	if oldUsername == newUsername {
		return fmt.Errorf("new username must be different from the old username")
	}

	// Start a database transaction to ensure atomicity
	tx, err := db.c.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback() // Rollback if an error occurs
			if rollbackErr != nil {
				// Rollbackerror + Original Error
				err = fmt.Errorf("failed to rollback transaction: %w, original error: %w", rollbackErr, err)
			}
		}
	}()

	// Update the user's username in the users table
	_, err = tx.Exec("UPDATE users SET username = ? WHERE username = ?", newUsername, oldUsername)
	if err != nil {
		return fmt.Errorf("failed to update username: %w", err)
	}

	// Update the username in all associated conversations (user1, user2)
	_, err = tx.Exec("UPDATE conversations SET user1 = ? WHERE user1 = ?", newUsername, oldUsername)
	if err != nil {
		return fmt.Errorf("failed to update username in conversations (user1): %w", err)
	}

	_, err = tx.Exec("UPDATE conversations SET user2 = ? WHERE user2 = ?", newUsername, oldUsername)
	if err != nil {
		return fmt.Errorf("failed to update username in conversations (user2): %w", err)
	}

	// Update the username in the messages table (sender)
	_, err = tx.Exec("UPDATE messages SET sender = ? WHERE sender = ?", newUsername, oldUsername)
	if err != nil {
		return fmt.Errorf("failed to update username in messages: %w", err)
	}

	// Update the username in message_status table
	_, err = tx.Exec("UPDATE message_status SET user_id = ? WHERE user_id = ?", newUsername, oldUsername)
	if err != nil {
		return fmt.Errorf("failed to update username in message_status: %w", err)
	}

	// Update the username in the group_members table for all the groups the user is a member of
	_, err = tx.Exec("UPDATE group_members SET membername = ? WHERE membername = ?", newUsername, oldUsername)
	if err != nil {
		return fmt.Errorf("failed to update username in group_members: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return nil if everything was successful
	return nil
}
