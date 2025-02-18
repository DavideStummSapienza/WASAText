package database

import (
	"errors"
	"fmt"
)

// ChangeGroupName changes the name of an existing group in the database.
//
// Parameters:
// - oldGroupName: The current name of the group.
// - newGroupName: The new name to be assigned to the group.
//
// Returns:
// - error: If an error occurs during the process (e.g., group does not exist, new name is already taken, etc.), an error is returned.
func (db *appdbimpl) ChangeGroupName(oldGroupName string, newGroupName string) error {
	// Check if the old group exists
	var groupExists int
	err := db.c.QueryRow(`
		SELECT COUNT(*) 
		FROM groups 
		WHERE groupname = ?;
	`, oldGroupName).Scan(&groupExists)

	if err != nil {
		return fmt.Errorf("failed to check if old group exists: %w", err)
	}

	// If the group does not exist, return an error
	if groupExists == 0 {
		return fmt.Errorf("group with name '%s' does not exist", oldGroupName)
	}

	// Check if the new groupname already exists as a username
	_, err = db.GetUser(newGroupName)

	if err == nil {
		return fmt.Errorf("group name already exists as a username")
	} else if !errors.Is(err, ErrUserNotFound) {
		return fmt.Errorf("database error while checking username: %w", err)
	}

	// Check if the new group name already exists
	var newGroupExists int
	err = db.c.QueryRow(`
		SELECT COUNT(*) 
		FROM groups 
		WHERE groupname = ?;
	`, newGroupName).Scan(&newGroupExists)

	if err != nil {
		return fmt.Errorf("failed to check if new group name exists: %w", err)
	}

	// If the new group name already exists, return an error
	if newGroupExists > 0 {
		return fmt.Errorf("group with name '%s' already exists", newGroupName)
	}

	// Now update the group name in the database
	_, err = db.c.Exec(`
		UPDATE groups 
		SET groupname = ? 
		WHERE groupname = ?;
	`, newGroupName, oldGroupName)

	if err != nil {
		return fmt.Errorf("failed to update group name: %w", err)
	}

	// Return nil if the group name was successfully changed
	return nil
}
