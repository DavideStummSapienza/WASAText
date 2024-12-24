package database

import (
	"fmt"
)

// ChangeProfilePicture updates the profile picture URL for a given user.
// The new profile picture URL is provided as `newProfilePhotoURL`.
// This function assumes the user has already been checked at the API level for existence.
func (db *appdbimpl) ChangeProfilePicture(username, newProfilePhotoURL string) error {
	// Update the user's profile picture URL in the database
	_, err := db.c.Exec("UPDATE users SET profile_photo_url = ? WHERE username = ?", newProfilePhotoURL, username)
	if err != nil {
		return fmt.Errorf("failed to update profile picture for user '%s': %w", username, err)
	}

	// Return nil if the update was successful
	return nil
}
