package database

import (
	"database/sql"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

// GetUser retrieves a user's details by their username.
// RETURNS: User struct if found, otherwise an error.
func (db *appdbimpl) GetUser(username string) (*User, error) {
	// SQL query to get the user by username
	query := "SELECT username, profile_photo_url, auth_token FROM users WHERE username = ?"

	// Struct to store the result
	var user User

	// Execute the query and scan the result into the user struct
	err := db.c.QueryRow(query, username).Scan(&user.Username, &user.ProfilePhotoURL, &user.AuthToken)
	if err != nil {
		// If no rows are found, return a custom error
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		// For other errors, return the original error
		return nil, err
	}

	// Return the found user
	return &user, nil
}
