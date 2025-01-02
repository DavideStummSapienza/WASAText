package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// AddComment adds a comment to a specific message in the conversation.
// It checks if the message exists and if the current user is a participant of the conversation.
// If the message is found and the user is allowed to comment, it inserts the comment into the `comments` table.
//
// Parameters:
// - messageTimestamp: The timestamp of the message being commented on.
// - partnerUsername: The username of the other participant in the conversation.
// - currentUser: The username of the user who is adding the comment.
// - content: The content of the comment (e.g., text or emoji).
//
// Returns:
// - error: If an error occurs during the process, such as a message not being found or a database failure, an error is returned.
func (db *appdbimpl) AddComment(messageTimestamp time.Time, partnerUsername string, currentUser string, content string) error {
	var messageID int

	// Check if the message exists and retrieve its ID by matching the timestamp and conversation participants.
	err := db.c.QueryRow(`
        SELECT m.id
        FROM messages m
        INNER JOIN conversations c ON m.id = c.message_id
        WHERE m.created_at = ? AND (c.from_user = ? OR c.to_user = ?)
    `, messageTimestamp, partnerUsername, partnerUsername).Scan(&messageID)

	// Handle case where no message was found
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("message not found: %w", err)
		}
		return fmt.Errorf("failed to fetch message ID: %w", err)
	}

	// Insert the comment into the comments table for the found message ID
	_, err = db.c.Exec(`
        INSERT INTO comments (reactor_username, message_id, content)
        VALUES (?, ?, ?)
    `, currentUser, messageID, content)

	// Handle any errors during the comment insertion
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	// Return nil if the comment was successfully added
	return nil
}
