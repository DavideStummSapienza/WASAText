package database

// Searches for a User in the Database
// Username can also be only partial
func (db *appdbimpl) SearchUser(partialUsername string) ([]string, error) {

	// Prepared Pattern for the search
	pattern := partialUsername + "%"
	rows, err := db.c.Query("SELECT username FROM users WHERE username LIKE ?", pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		usernames = append(usernames, username)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return usernames, nil
}
