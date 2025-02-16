package database

import "fmt"

// ShowConversation retrieves all messages for a specific conversation along with their metadata.
// It includes details like message content, sender information, timestamps, delivery/read status,
// and any reactions associated with each message. The function organizes messages in reverse
// chronological order, ensuring the newest messages are displayed first.
func (db *appdbimpl) ShowConversation(username, conversationPartnerName string) ([]ConversationDetail, error) {
	var conversation []ConversationDetail

	rows, err := db.c.Query(`
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
    WHERE m.conversation_id IN (
        SELECT id
        FROM conversations
        WHERE 
            (user1 = ? AND user2 = ?) 
            OR (user1 = ? AND user2 = ?) 
            OR groupname = ?
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

		if err := rows.Scan(&msg.MessageID, &msg.Content, &msg.Sender, &msg.IsPhoto, &msg.IsForwarded, &msg.Timestamp, &msg.FullyReceived, &msg.FullyRead); err != nil {
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
func (db *appdbimpl) getReactionsForMessage(messageID int) ([]Reaction, error) {
	var reactions []Reaction

	// Query to retrieve both the reactor's username and content of the reaction
	rows, err := db.c.Query(`
        SELECT reactor_username, content
        FROM comments
        WHERE message_id = ?`, messageID)
	if err != nil {
		return nil, fmt.Errorf("error querying reactions for message '%d': %w", messageID, err)
	}
	defer rows.Close()

	// Collect all reactions from the result set
	for rows.Next() {
		var reaction Reaction
		if err := rows.Scan(&reaction.Reactor, &reaction.Content); err != nil {
			return nil, fmt.Errorf("error scanning reaction: %w", err)
		}
		reactions = append(reactions, reaction)
	}

	// Check for errors while iterating through rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return reactions, nil
}
