package database

import "fmt"

// DeleteComment removes a comment from a specific message in the database.
// It checks if the message exists and if the comment belongs to the specified reactor.
// If the comment is found and the reactor is authorized to delete it, the comment is removed from the `comments` table.
//
// Parameters:
// - messageID: The ID of the message the comment is associated with.
// - reactorUsername: The username of the user who posted the comment that is to be deleted.
//
// Returns:
// - error: If an error occurs during the deletion, such as the comment not being found or a database failure, an error is returned.
func (db *appdbimpl) DeleteComment(messageID int, reactorUsername string) error {
	// Execute the DELETE statement to remove the comment where the message_id and reactor_id match
	_, err := db.c.Exec(`
		DELETE FROM comments
		WHERE message_id = ? AND reactor_username = ?;
	`, messageID, reactorUsername)

	// Handle any errors during the deletion
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	// Return nil if the comment was successfully deleted
	return nil
}
