package database

import "fmt"

// MarkAllMessagesAsReceived marks all messages in a conversation (1:1 or group) as received for a user.
func (db *appdbimpl) MarkAllMessagesAsReceived(conversationPartner string, username string) error {
	_, err := db.c.Exec(`
		UPDATE message_status
		SET received = TRUE
		WHERE message_id IN (
			SELECT m.id
			FROM messages m
			JOIN conversations c ON m.id = c.message_id
			WHERE 
				(c.to_user = ? AND c.from_user = ?) OR
				(c.to_group = ? AND ? IN (SELECT user FROM group_members WHERE groupname = c.to_group))
		) AND user_id = ?`,
		username, conversationPartner, 
		conversationPartner, username, 
		username)

	if err != nil {
		return fmt.Errorf("error marking messages as received: %w", err)
	}
	return nil
}
