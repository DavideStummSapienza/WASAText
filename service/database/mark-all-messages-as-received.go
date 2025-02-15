package database

import (
	"fmt"
	"log"
)

// MarkAllMessagesAsReceived updates all messages in a conversation as received for a user.
func (db *appdbimpl) MarkAllMessagesAsReceived(username string, partnerUsername string) error {
	_, err := db.c.Exec(`
		UPDATE message_status
		SET received = TRUE
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
		partnerUsername, username, // Reversed
		username, // Groupchat
		username, // User whose status is updated
	)

	if err != nil {
		log.Printf("Error marking messages as received for user %s and partner %s: %v", username, partnerUsername, err)
		return fmt.Errorf("error marking messages as received: %w", err)
	}
	return nil
}
