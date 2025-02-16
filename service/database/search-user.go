package database

import "database/sql"

// Searches for a User in the Database
// Username can also be only partial
// RETURNS: all usernames or partialusernames, if found otherwise an error will be returned
func (db *appdbimpl) SearchUser(partialUsername, currentUser string) ([]User, error) {
	var query string
	var rows *sql.Rows
	var err error

	// If partialUsername is empty, retrieve all users from the database
	if partialUsername == "" {
		query = "SELECT username, profile_photo_url FROM users WHERE username != ?"
		rows, err = db.c.Query(query, currentUser)
	} else {
		// If partialUsername is not empty, use LIKE for searching
		pattern := partialUsername + "%"
		query = "SELECT username, profile_photo_url FROM users WHERE username LIKE ? AND username != ?"
		rows, err = db.c.Query(query, pattern, currentUser)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through all the rows returned from the query
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Username, &user.ProfilePhotoURL); err != nil {
			return nil, err
		}
		users = append(users, user) // Appending the user to the slice
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
