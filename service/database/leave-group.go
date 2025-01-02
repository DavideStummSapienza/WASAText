package database

import (
	"fmt"
)

// LeaveGroup removes a user from an existing group.
//
// Parameters:
// - groupName: The name of the group that the user wants to leave.
// - currentUser: The username of the user who wants to leave the group.
//
// Returns:
// - error: If an error occurs during the process (e.g., group does not exist, user is not a member, or failure to remove user), an error is returned.
func (db *appdbimpl) LeaveGroup(groupName string, currentUser string) error {
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

	// Check if the user is a member of the group
	var userIsMember int
	err = db.c.QueryRow(`
		SELECT COUNT(*) 
		FROM group_members 
		WHERE groupname = ? AND membername = ?;
	`, groupName, currentUser).Scan(&userIsMember)

	if err != nil {
		return fmt.Errorf("failed to check if user is a member of the group: %w", err)
	}

	// If the user is not a member of the group, return an error
	if userIsMember == 0 {
		return fmt.Errorf("user '%s' is not a member of group '%s'", currentUser, groupName)
	}

	// Remove the user from the group
	_, err = db.c.Exec(`
		DELETE FROM group_members 
		WHERE groupname = ? AND membername = ?;
	`, groupName, currentUser)

	if err != nil {
		return fmt.Errorf("failed to remove user from group: %w", err)
	}

	// Return nil if the user was successfully removed from the group
	return nil
}
