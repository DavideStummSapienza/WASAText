package database

import (
	"fmt"
)

// DeleteComment removes a user's comment from a message.
//
// Parameters:
// - messageID: The ID of the message the comment is associated with.
// - reactorUsername: The username of the user who posted the comment.
//
// Returns:
// - error: If the comment is not found or a database error occurs.
func (db *appdbimpl) DeleteComment(messageID int, reactorUsername string) error {
	// Execute the DELETE statement
	res, err := db.c.Exec(`
		DELETE FROM comments
		WHERE message_id = ? AND reactor_username = ?;
	`, messageID, reactorUsername)

	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment not found or user not authorized")
	}

	return nil
}
