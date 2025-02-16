package database

import (
	"database/sql"
	"errors"
)

var ErrGroupNotFound = errors.New("group not found")

// GetGroupByName retrieves a group by its name from the database.
// Returns: Group struct if found, otherwise an error.
func (db *appdbimpl) GetGroupByName(groupName string) (*Group, error) {
	query := "SELECT groupname, group_photo_url FROM groups WHERE groupname = ?"
	var group Group

	// Execute query and scan result into the group struct
	err := db.c.QueryRow(query, groupName).Scan(&group.Groupname, &group.GroupPhotoUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}
	return &group, nil
}
