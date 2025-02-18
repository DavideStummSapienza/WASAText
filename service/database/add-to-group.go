package database

import (
	"errors"
	"fmt"
	"log"
)

// AddToGroup adds multiple users to a group.
// If the group does not exist, it will be created, and the provided users (including the current user) will be added to the group.
//
// Parameters:
// - groupname: The name of the group to which users should be added or created.
// - usernames: A list of usernames to be added to the group.
// - currentUser: The username of the user who is making the request to add users to the group.
//
// Returns:
// - error: If an error occurs during the process, such as a database failure or permission issue, an error is returned.
func (db *appdbimpl) AddToGroup(groupname string, usernames []string, currentUser string) error {

	// Check if the groupname already exists as a username
	_, err := db.GetUser(groupname)
	if err == nil {
		return fmt.Errorf("group name already exists as a username")
	} else if !errors.Is(err, ErrUserNotFound) {
		return fmt.Errorf("database error while checking username: %w", err)
	}

	// Check if the group exists
	var groupCount int
	err = db.c.QueryRow(`
		SELECT COUNT(*) 
		FROM groups 
		WHERE groupname = ?;
	`, groupname).Scan(&groupCount)

	if err != nil {
		return fmt.Errorf("failed to check if group exists: %w", err)
	}

	// If the group doesn't exist, create it
	if groupCount == 0 {

		log.Print("Group created")

		_, err := db.c.Exec(`
			INSERT INTO groups (groupname)
			VALUES (?);
		`, groupname)

		if err != nil {
			return fmt.Errorf("failed to create group: %w", err)
		}

		// Create the conversation for the new group

		_, err = db.c.Exec(`
			INSERT INTO conversations (groupname)
			VALUES (?);
		`, groupname)

		if err != nil {
			return fmt.Errorf("failed to create conversation for new group: %w", err)
		}

		log.Print("Conversation for group created")

		// Add the current user to the group
		_, err = db.c.Exec(`
			INSERT INTO group_members (groupname, membername)
			VALUES (?, ?);
		`, groupname, currentUser)

		if err != nil {
			return fmt.Errorf("failed to add current user to group: %w", err)
		}

		log.Printf("current user added to new group %s %s", groupname, currentUser)
	}

	// Add the other users to the group (if not already members)
	for _, username := range usernames {
		// Check if the user is already a member of the group
		var userAlreadyMember int
		err = db.c.QueryRow(`
			SELECT COUNT(*) 
			FROM group_members 
			WHERE groupname = ? AND membername = ?;
		`, groupname, username).Scan(&userAlreadyMember)

		if err != nil {
			return fmt.Errorf("failed to check if user is already a member: %w", err)
		}

		if userAlreadyMember > 0 {
			continue // Skip adding the user if they are already a member
		}

		// Add the user to the group
		_, err = db.c.Exec(`
			INSERT INTO group_members (groupname, membername)
			VALUES (?, ?);
		`, groupname, username)

		if err != nil {
			return fmt.Errorf("failed to add user %s to group: %w", username, err)
		}
	}

	// Return nil if all users were successfully added
	return nil
}
