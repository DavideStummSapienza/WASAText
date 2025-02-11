package database

import (
	"fmt"
	"log"
)

// MarkAllMessagesAsRead updates all messages in a conversation as read for a user.
func (db *appdbimpl) MarkAllMessagesAsRead(username string, partnerUsername string) error {


	_, err := db.c.Exec(`
		UPDATE message_status
		SET read = TRUE
		WHERE message_id IN (
			SELECT m.id
			FROM messages m
			JOIN conversations c ON m.id = c.message_id
			WHERE 
				(c.to_user = ? AND c.from_user = ?) OR
				(c.to_group = ? AND ? IN (SELECT membername FROM group_members WHERE groupname = c.to_group))
		) AND user_id = ?`,
		username, partnerUsername,
		partnerUsername, username,
		username)

	if err != nil {
		log.Printf("Error marking messages as read for user %s and partner %s: %v", username, partnerUsername, err)
		return fmt.Errorf("error marking messages as read: %w", err)
	}
	return nil
}
