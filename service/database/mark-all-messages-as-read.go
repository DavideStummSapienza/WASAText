package database

import "fmt"

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
				(c.to_user = ? AND c.from_user = ?) OR
				(c.to_group = ? AND ? IN (SELECT user FROM group_members WHERE groupname = c.to_group))
		) AND user_id = ?`,
		username, partnerUsername,
		partnerUsername, username,
		username)

	if err != nil {
		return fmt.Errorf("error marking messages as read: %w", err)
	}
	return nil
}
