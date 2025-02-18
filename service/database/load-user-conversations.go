package database

import "time"

// LoadUserConversations fetches a list of conversation previews for the given user.
func (db *appdbimpl) LoadUserConversations(username string) ([]ConversationPreview, error) {
	// SQL query to retrieve the required data for the conversation preview
	query := `
	SELECT 
		CASE 
			WHEN c.groupname IS NOT NULL THEN g.groupname
			ELSE 
				CASE 
					WHEN c.user1 = ? THEN c.user2
					ELSE c.user1 
				END
		END AS name,
		CASE 
			WHEN c.groupname IS NOT NULL THEN g.group_photo_url
			ELSE 
				CASE 
					WHEN c.user1 = ? THEN u2.profile_photo_url 
					ELSE u1.profile_photo_url 
				END
		END AS photo_url,
		CASE
			WHEN m.is_photo THEN 'Photo'
			WHEN m.content IS NULL THEN 'No messages yet'
        	ELSE m.content
		END AS last_message,
		m.created_at AS last_message_time,
		CASE 
    		WHEN c.groupname IS NOT NULL THEN true
    		ELSE false
		END AS is_group
	FROM conversations c
	LEFT JOIN messages m ON m.id = (
		SELECT id FROM messages 
		WHERE conversation_id = c.id 
		ORDER BY created_at DESC 
		LIMIT 1
	)
	LEFT JOIN users u1 ON u1.username = c.user1
	LEFT JOIN users u2 ON u2.username = c.user2
	LEFT JOIN groups g ON g.groupname = c.groupname
	LEFT JOIN group_members gm ON gm.groupname = c.groupname
	WHERE c.user1 = ? OR c.user2 = ? OR gm.membername = ?
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
		err := rows.Scan(&preview.Name, &preview.PhotoURL, &preview.LastMessage, &preview.LastMessageTime, &preview.IsGroup)
		if err != nil {
			return nil, err
		}

		if !preview.LastMessageTime.Valid {
			preview.LastMessageTime.Time = time.Time{} // Default "zero time"
		}

		previews = append(previews, preview)
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return previews, nil
}
