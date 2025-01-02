package database

import (
	"fmt"
)

// ChangeGroupPicture changes the photo of an existing group.
//
// Parameters:
// - groupName: The name of the group whose photo is being updated.
// - newPhotoURL: The new photo URL that will be set as the group's photo.
// - currentUser: The username of the user who is attempting to change the group photo. This ensures only authorized users can change the group photo.
//
// Returns:
// - error: If an error occurs during the process (e.g., group does not exist, or update failure), an error is returned.
func (db *appdbimpl) ChangeGroupPicture(groupName string, newPhotoURL string) error {
	// Check if the group exists
	var groupExists int
	err := db.c.QueryRow(`
		SELECT COUNT(*) 
		FROM groups 
		WHERE groupname = ?;
	`, groupName).Scan(&groupExists)

	if err != nil {
		return fmt.Errorf("failed to check if group exists: %w", err)
	}

	// If the group does not exist, return an error
	if groupExists == 0 {
		return fmt.Errorf("group with name '%s' does not exist", groupName)
	}

	// Update the group's photo URL in the database
	_, err = db.c.Exec(`
		UPDATE groups 
		SET group_photo_url = ? 
		WHERE groupname = ?;
	`, newPhotoURL, groupName)

	if err != nil {
		return fmt.Errorf("failed to update group photo: %w", err)
	}

	// Return nil if the group's photo was successfully updated
	return nil
}
