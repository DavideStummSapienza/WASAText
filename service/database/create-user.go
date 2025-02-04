package database

// CreateUser inserts a new user into the database.
func (db *appdbimpl) CreateUser(username string, profilePhotoURL string, authToken int) error {
	_, err := db.c.Exec("INSERT INTO users (username, profile_photo_url, auth_token) VALUES (?, ?, ?)", username, profilePhotoURL, authToken)
	return err
}
