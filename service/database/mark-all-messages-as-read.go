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
			JOIN conversations c ON m.conversation_id = c.id
			WHERE 
				(c.user1 = ? AND c.user2 = ?) OR
				(c.user2 = ? AND c.user1 = ?) OR
				(c.groupname IS NOT NULL AND c.groupname IN (
					SELECT groupname FROM group_members WHERE membername = ?
				))
		) 
		AND user_id = ?`,
		username, partnerUsername, // Privatechat
		partnerUsername, username, // Reverse
		username, // Groupchat
		username, // User whose status is updated
	)

	if err != nil {
		log.Printf("Error marking messages as read for user %s and partner %s: %v", username, partnerUsername, err)
		return fmt.Errorf("error marking messages as read: %w", err)
	}
	return nil
}
