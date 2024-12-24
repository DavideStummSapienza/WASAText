package database

// LoadUserConversations fetches a list of conversation previews for the given user.
func (db *appdbimpl) LoadUserConversations(username string) ([]ConversationPreview, error) {
	// SQL query to retrieve the required data for the conversation preview
	query := `
	SELECT 
		CASE
			WHEN c.to_group IS NOT NULL THEN g.groupname
			ELSE u.username
		END AS name,
		CASE
			WHEN c.to_group IS NOT NULL THEN g.group_photo_url
			ELSE u.profile_photo_url
		END AS photo_url,
		CASE
			WHEN m.is_photo THEN '📷 Photo'
			ELSE m.content
		END AS last_message,
		m.created_at AS last_message_time
	FROM (
		SELECT 
			c.id AS conversation_id,
			MAX(m.created_at) AS latest_message_time
		FROM 
			conversations c
		JOIN 
			messages m ON c.message_id = m.id
		WHERE 
			c.from_user = ? OR c.to_user = ?
		GROUP BY 
			c.id
	) latest_conversations
	JOIN conversations c ON c.id = latest_conversations.conversation_id
	LEFT JOIN users u ON (u.username = c.to_user OR u.username = c.from_user)
	LEFT JOIN groups g ON g.groupname = c.to_group
	JOIN messages m ON m.created_at = latest_conversations.latest_message_time
	ORDER BY 
		latest_conversations.latest_message_time DESC;
	`

	// Execute the query
	rows, err := db.c.Query(query, username, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Prepare the list of conversation previews
	var previews []ConversationPreview
	for rows.Next() {
		var preview ConversationPreview
		err := rows.Scan(&preview.Name, &preview.PhotoURL, &preview.LastMessage, &preview.LastMessageTime)
		if err != nil {
			return nil, err
		}
		previews = append(previews, preview)
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return previews, nil
}
