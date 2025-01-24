package database

import (
	"database/sql"
	"errors"
)

// GetUsernameByToken prüft, ob das Token existiert, und gibt den zugehörigen Benutzernamen zurück
func (db *appdbimpl) GetUsernameByToken(token string) (string, error) {
    var username string
    query := `SELECT username FROM users WHERE auth_token = ?`
    err := db.c.QueryRow(query, token).Scan(&username)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            // Kein Benutzer mit diesem Token gefunden
            return "", nil
        }
        // Anderer Datenbankfehler
        return "", err
    }
    // Benutzername gefunden
    return username, nil
}
