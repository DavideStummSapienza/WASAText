package database

import (
	"database/sql"
	"errors"
)

// GetUsernameByToken checks if the token exists and returns the associated username.
// If the token does not exist, it returns an empty string and no error.
// If there is a database error, it returns the error.
func (db *appdbimpl) GetUsernameByToken(token string) (string, error) {
	var username string
	// SQL query to find the username associated with the provided token.
	query := `SELECT username FROM users WHERE auth_token = ?`
	err := db.c.QueryRow(query, token).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No user found with this token; return an empty username and no error.
			return "", nil
		}
		// Some other database error occurred; return the error.
		return "", err
	}
	// Username successfully found; return it.
	return username, nil
}
