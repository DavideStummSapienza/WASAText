package database

import (
	"database/sql"
	"fmt"
)

// ShowConversation retrieves all messages and associated details for a given conversation.
func (db *appdbimpl) ShowConversation(username, conversationID string) ([]ConversationDetail, error) {
	// Initialize a slice to hold the detailed conversation data
	var conversation []ConversationDetail

	// Query to get all messages for the conversation
	// We assume 'conversationID' can refer to either a 1:1 conversation or a group conversation
	rows, err := db.c.Query(`
		SELECT 
			m.id, 
			m.content, 
			m.is_photo, 
			m.photo_url, 
			m.created_at, 
			ms.user_id, 
			ms.received, 
			ms.read, 
			c.content as reaction_content
		FROM messages m
		LEFT JOIN message_status ms ON m.id = ms.message_id
		LEFT JOIN comments c ON m.id = c.message_id
		WHERE m.id IN (
			SELECT message_id
			FROM conversations
			WHERE (from_user = ? OR to_user = ? OR to_group = ?)
		)
		ORDER BY m.created_at DESC`, username, username, conversationID)
	if err != nil {
		return nil, fmt.Errorf("error querying messages for conversation '%s': %w", conversationID, err)
	}
	defer rows.Close()

	// Loop over all rows in the result set and construct the conversation details
	for rows.Next() {
		var msg ConversationDetail
		var sender string
		var received, read bool
		var reactionContent sql.NullString // To handle potential NULL reactions

		// Scan row values into the structure
		if err := rows.Scan(&msg.MessageID, &msg.Content, &msg.IsPhoto, &msg.PhotoURL, &msg.Timestamp, &sender, &received, &read, &reactionContent); err != nil {
			return nil, fmt.Errorf("error scanning message row: %w", err)
		}

		// Set the sender of the message
		msg.Sender = sender

		// Determine message status based on received and read flags
		if received && read {
			msg.Status = "Read"
		} else if received {
			msg.Status = "Received"
		} else {
			msg.Status = "Sent"
		}

		// Collect reactions (comments) for the message
		if reactionContent.Valid {
			msg.Reactions = append(msg.Reactions, reactionContent.String)
		}

		// Append the message to the conversation slice
		conversation = append(conversation, msg)
	}

	// Check for errors while iterating through the rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error processing rows: %w", err)
	}

	// Return the detailed conversation
	return conversation, nil
}
