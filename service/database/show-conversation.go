package database

import "fmt"

// ShowConversation retrieves all messages for a specific conversation along with their metadata.
// It includes details like message content, sender information, timestamps, delivery/read status,
// and any reactions associated with each message. The function organizes messages in reverse
// chronological order, ensuring the newest messages are displayed first.
func (db *appdbimpl) ShowConversation(username, conversationPartnerName string) ([]ConversationDetail, error) {
	var conversation []ConversationDetail

	// 1. Retrieve all messages for the conversation
	rows, err := db.c.Query(`
    SELECT 
        m.id, 
        m.content, 
        m.is_photo, 
        m.photo_url, 
        m.created_at,
		c.from_user, 
        ms.user_id, 
        ms.received, 
        ms.read
    FROM messages m
    LEFT JOIN message_status ms ON m.id = ms.message_id
	LEFT JOIN conversations c ON m.id = c.message_id
    WHERE m.id IN (
        SELECT message_id
        FROM conversations
        WHERE 
            (from_user = ? AND to_user = ?) 
            OR (from_user = ? AND to_user = ?) 
            OR to_group = ?
    )
    ORDER BY m.created_at DESC`, 
    username, conversationPartnerName, conversationPartnerName, username, conversationPartnerName)
	
	if err != nil {
		return nil, fmt.Errorf("error querying messages for conversation '%s': %w", conversationPartnerName, err)
	}
	defer rows.Close()

	// Populate the conversation details
	for rows.Next() {
		var msg ConversationDetail

		if err := rows.Scan(&msg.MessageID, &msg.Content, &msg.IsPhoto, &msg.PhotoURL, &msg.Timestamp, &msg.Sender, &msg.Receiver, &msg.Received, &msg.Read); err != nil {
			return nil, fmt.Errorf("error scanning message row: %w", err)
		}

		// 2. Retrieve reactions for the current message
		reactions, err := db.getReactionsForMessage(msg.MessageID)
		if err != nil {
			return nil, fmt.Errorf("error retrieving reactions for message '%d': %w", msg.MessageID, err)
		}
		msg.Reactions = reactions

		// Add the message to the conversation list
		conversation = append(conversation, msg)
	}

	// Check for errors while iterating through rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error processing rows: %w", err)
	}

	return conversation, nil
}

// Helper function: Retrieve all reactions for a given message
func (db *appdbimpl) getReactionsForMessage(messageID int) ([]string, error) {
	var reactions []string

	// Query to retrieve reactions
	rows, err := db.c.Query(`
        SELECT content
        FROM comments
        WHERE message_id = ?`, messageID)
	if err != nil {
		return nil, fmt.Errorf("error querying reactions for message '%d': %w", messageID, err)
	}
	defer rows.Close()

	// Collect all reactions from the result set
	for rows.Next() {
		var reaction string
		if err := rows.Scan(&reaction); err != nil {
			return nil, fmt.Errorf("error scanning reaction: %w", err)
		}
		reactions = append(reactions, reaction)
	}

	return reactions, nil
}
