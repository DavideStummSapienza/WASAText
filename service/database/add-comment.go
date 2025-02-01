package database

import (
	"fmt"
)

// AddComment inserts a comment into the `comments` table for a given message.
//
// Parameters:
// - messageID: The ID of the message to which the comment is added.
// - currentUser: The username of the user who is making the comment.
// - content: The comment text or emoji.
//
// Returns:
// - error: An error if the database insertion fails.
func (db *appdbimpl) AddComment(messageID int, currentUser string, content string) error {
	_, err := db.c.Exec(`
        INSERT INTO comments (reactor_username, message_id, content)
        VALUES (?, ?, ?)`, currentUser, messageID, content)

	if err != nil {
		return fmt.Errorf("failed to insert comment: %w", err)
	}

	return nil
}
