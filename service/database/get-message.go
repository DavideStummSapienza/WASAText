package database

import "fmt"

// GetMessage retrieves a specific message by its ID or the latest message in a conversation.
//
// If a message ID is provided, the function fetches the corresponding message from the database.
// If no message ID is provided, it retrieves the most recent message exchanged between `username` and `partner`.
//
// The function also fetches metadata such as whether the message was fully received or read,
// along with any reactions associated with it.
//
// Parameters:
// - messageID: (optional) The ID of the message to fetch. If nil, the latest message is retrieved.
// - username: The username of the requesting user.
// - partner: The username of the conversation partner (or group name).
//
// Returns:
// - A pointer to a ConversationDetail struct containing message details.
// - An error if the message retrieval fails.
func (db *appdbimpl) GetMessage(messageID *int) (*ConversationDetail, error) {
	var msg ConversationDetail
	var err error

	// If a specific message ID is provided, retrieve that message
	err = db.c.QueryRow(`
    SELECT 
        m.id, 
        m.content,
        m.sender, 
        m.is_photo, 
        m.is_forwarded, 
        m.created_at, 
        (SELECT COUNT(*) FROM message_status WHERE message_id = m.id AND received = FALSE) = 0 AS fully_received,
        (SELECT COUNT(*) FROM message_status WHERE message_id = m.id AND read = FALSE) = 0 AS fully_read
    FROM messages m
    WHERE m.id = ?`, *messageID).Scan(&msg.MessageID, &msg.Content, &msg.Sender, &msg.IsPhoto, &msg.IsForwarded, &msg.Timestamp, &msg.FullyReceived, &msg.FullyRead)

	// Handle any errors during the database query
	if err != nil {
		return nil, fmt.Errorf("error retrieving message: %w", err)
	}

	// Retrieve reactions for the message
	msg.Reactions, err = db.getReactionsForMessage(msg.MessageID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving reactions: %w", err)
	}

	return &msg, nil
}
