/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	//Securityralated Functions
	GetUsernameByToken(token string) (string, error)

	// Userrelated Functions
	GetUser(username string) (*User, error)
	CreateUser(username string, profilePhotoURL string, authToken int) error
	SearchUser(partialUsername string) ([]User, error)
	LoadUserConversations(username string) ([]ConversationPreview, error)
	ChangeUsername(oldUsername, newUsername string) error
	ChangeProfilePicture(username, newProfilePhotoURL string) error
	GetGroupByName(groupName string) (*Group, error)

	// Conversation Functions
	ShowConversation(username, conversationPartnerName string) ([]ConversationDetail, error)
	SendMessage(msg NewMessage) (int, error)
	DeleteMessage(currentUser string, messageID int) error
	MarkAllMessagesAsReceived(partnerUsername string, username string) error
	MarkAllMessagesAsRead(username string, partnerUsername string) error
	GetMessage(messageID *int) (*ConversationDetail, error)

	// Comment Functions
	AddComment(messageID int, currentUser string, content string) error
	DeleteComment(messageID int, reactorUsername string) error

	// Group Functions
	AddToGroup(groupname string, usernames []string, currentUser string) error
	ChangeGroupPicture(groupName string, newPhotoURL string) error
	ChangeGroupName(oldGroupName string, newGroupName string) error
	LeaveGroup(groupName string, currentUser string) error

	GetName() (string, error)
	SetName(name string) error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building an AppDatabase")
	}

	// List of Tablenames and their SQL-CREATE-Commands
	tables := map[string]string{
		"users": `
			CREATE TABLE IF NOT EXISTS users (
				username TEXT UNIQUE NOT NULL PRIMARY KEY,
				profile_photo_url TEXT,
				auth_token INTEGER UNIQUE NOT NULL
			);
		`,
		"conversations": `
			CREATE TABLE IF NOT EXISTS conversations (
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				user1 TEXT REFERENCES users(username) ON DELETE CASCADE,
				user2 TEXT REFERENCES users(username) ON DELETE CASCADE,
				groupname TEXT REFERENCES groups(groupname) ON DELETE CASCADE
			);
		`,
		"messages": `
			CREATE TABLE IF NOT EXISTS messages (
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				content TEXT,
				sender TEXT REFERENCES users(username) ON DELETE CASCADE,
				is_photo BOOLEAN DEFAULT FALSE,
				is_forwarded BOOLEAN DEFAULT FALSE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				conversation_id INTEGER NOT NULL REFERENCES conversations(id) ON DELETE CASCADE
			);
		`,
		"message_status": `
			CREATE TABLE IF NOT EXISTS message_status (
				message_id INTEGER NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
				user_id TEXT NOT NULL REFERENCES users(username) ON DELETE CASCADE,
				received BOOLEAN DEFAULT FALSE,
				read BOOLEAN DEFAULT FALSE,
				PRIMARY KEY (message_id, user_id)
			);
		`,
		"comments": `
			CREATE TABLE IF NOT EXISTS comments (
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				reactor_username TEXT NOT NULL REFERENCES users(username) ON DELETE CASCADE,
				message_id INTEGER NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
				content TEXT,
				UNIQUE (reactor_username, message_id)
			);
		`,
		"groups": `
			CREATE TABLE IF NOT EXISTS groups (
				groupname TEXT NOT NULL PRIMARY KEY,
				group_photo_url TEXT
			);
		`,
		"group_members": `
			CREATE TABLE IF NOT EXISTS group_members (
				groupname TEXT NOT NULL REFERENCES groups(groupname) ON DELETE CASCADE,
				membername TEXT NOT NULL REFERENCES users(username) ON DELETE CASCADE
			);
		`,
	}

	// Check if table exists. If not, the database will be created

	for tableName, createStmt := range tables {
		_, err := db.Exec(createStmt)
		if err != nil {
			fmt.Printf("SQL Error: Failed to create table %s. Statement: %s, Error: %v\n", tableName, createStmt, err)
			return nil, fmt.Errorf("error creating table %s: %w", tableName, err)
		}
	}

	/*
		// Check if table exists. If not, the database is empty, and we need to create the structure

		var tableName string
		err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
		if errors.Is(err, sql.ErrNoRows) {
			sqlStmt := `CREATE TABLE example_table (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name TEXT);`
			_, err = db.Exec(sqlStmt)
			if err != nil {
				return nil, fmt.Errorf("error creating database structure: %w", err)
			}
		}

	*/

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
