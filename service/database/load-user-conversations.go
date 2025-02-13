package database

// LoadUserConversations fetches a list of conversation previews for the given user.
func (db *appdbimpl) LoadUserConversations(username string) ([]ConversationPreview, error) {
	// SQL query to retrieve the required data for the conversation preview
	query := `
	SELECT 
		CASE 
			WHEN c.to_group IS NOT NULL THEN g.groupname
			ELSE 
				CASE 
					WHEN c.from_user = ? THEN c.to_user 
					ELSE c.from_user 
				END
		END AS name,
		CASE 
			WHEN c.to_group IS NOT NULL THEN g.group_photo_url
			ELSE 
				CASE 
					WHEN c.from_user = ? THEN u2.profile_photo_url 
					ELSE u1.profile_photo_url 
				END
		END AS photo_url,
		CASE
			WHEN m.is_photo THEN 'Photo'
			ELSE m.content
		END AS last_message,
		m.created_at AS last_message_time
	FROM conversations c
	JOIN messages m ON m.id = (
		SELECT id FROM messages 
		WHERE conversation_id = c.id 
		ORDER BY created_at DESC 
		LIMIT 1
	)
	LEFT JOIN users u1 ON u1.username = c.from_user
	LEFT JOIN users u2 ON u2.username = c.to_user
	LEFT JOIN groups g ON g.groupname = c.to_group
	WHERE c.from_user = ? OR c.to_user = ? OR c.to_group = ?
	ORDER BY m.created_at DESC;
	`

	// Execute the query
	rows, err := db.c.Query(query, username, username, username, username, username)
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
