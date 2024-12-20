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
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// List of Tablenames and their SQL-CREATE-Commands
	tables := map[string]string{
		"users": `
			CREATE TABLE IF NOT EXISTS users (
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				username TEXT UNIQUE NOT NULL,
				name TEXT,
				profile_photo_url TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
		`,
		"conversations": `
			CREATE TABLE IF NOT EXISTS conversations (
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				name TEXT,
				is_group BOOLEAN DEFAULT FALSE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
		`,
		"conversation_members": `
			CREATE TABLE IF NOT EXISTS conversation_members (
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				conversation_id INTEGER NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
				user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				UNIQUE (conversation_id, user_id)
			);
		`,
		"messages": `
			CREATE TABLE IF NOT EXISTS messages (
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				conversation_id INTEGER NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
				sender_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				content TEXT,
				photo_url TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
		`,
		"message_status": `
			CREATE TABLE IF NOT EXISTS message_status (
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				message_id INTEGER NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
				user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				status TEXT CHECK(status IN ('sent', 'received', 'read')) DEFAULT 'sent',
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				UNIQUE (message_id, user_id)
			);
		`,
		"reactions": `
			CREATE TABLE IF NOT EXISTS reactions (
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				message_id INTEGER NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
				user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				reaction_type TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
		`,
	}

	// Check if table exists. If not, the database will be created

	for tableName, createStmt := range tables {
		_, err := db.Exec(createStmt)
		if err != nil {
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
