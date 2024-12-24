package database

// Searches for a User in the Database
// Username can also be only partial
// RETRUNS: all usernames or partialusernames, if found otherwise an error will be returned
func (db *appdbimpl) SearchUser(partialUsername string) ([]string, error) {

	// Prepared Pattern for the search
	pattern := partialUsername + "%"

	// First row will be checked, if the slice is part of the username
	rows, err := db.c.Query("SELECT username FROM users WHERE username LIKE ?", pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Check runs for all the rows
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
