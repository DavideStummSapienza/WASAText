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
			tx.Rollback() // Rollback if an error occurs
		}
	}()

	// Update the user's username in the users table
	_, err = tx.Exec("UPDATE users SET username = ? WHERE username = ?", newUsername, oldUsername)
	if err != nil {
		return fmt.Errorf("failed to update username: %w", err)
	}

	// Update the username in all associated conversations
	_, err = tx.Exec("UPDATE conversations SET from_user = ? WHERE from_user = ?", newUsername, oldUsername)
	if err != nil {
		return fmt.Errorf("failed to update username in conversations: %w", err)
	}

	// Update the username in message status table
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
