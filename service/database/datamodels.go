package database

import "time"

// ConversationPreview represents the preview information for a conversation.
type ConversationPreview struct {
	Name            string    `json:"name"`              // Username or Group Name
	PhotoURL        string    `json:"photo_url"`         // Profile Photo URL
	LastMessage     string    `json:"last_message"`      // Snippet or Photo Icon
	LastMessageTime time.Time `json:"last_message_time"` // Timestamp of Last Message
}

// ConversationDetail represents the detailed information for a message in a conversation.
type ConversationDetail struct {
	MessageID 		int       `json:"message_id"` 	 	// ID of the message
	Content   		string    `json:"content"`    	 	// Content of the message (text or photo URL)
	IsPhoto   		bool      `json:"is_photo"`   	 	// Whether the message is a photo message
	PhotoURL 		string    `json:"photo_url"`  	 	// URL of the photo (if IsPhoto is true)
	IsForwarded		bool	  `json:"is_forwarded"`  	// Whether the message is a photo message
	Timestamp 		time.Time `json:"timestamp"`  	 	// Timestamp of when the message was created
	Sender    		string    `json:"sender"`     	 	// Sender of the message
	Receiver  		string	  `json:"receiver"`   	 	// Receiver of the message
	FullyReceived   bool 	  `json:"fully_received"`   // Received-Status of the message
	FullyRead       bool 	  `json:"fully_read"`       // Read-Status of the message
	Reactions 	    []Reaction  `json:"reactions"`  		// List of user reactions (comments)
}

// Reaction represents the Reaction on a message in a conversation.
type Reaction struct {
	Reactor string `json:"reactor"` // username of the Reactor
	Content string `json:"content"` // Content of the Reaction (e.g. Emoji)
}


// NewMessage is used for the Parameters of the SendMessage Function
type NewMessage struct {
	FromUser    string
    ToUser      string
    Content     string
    IsPhoto     bool
    PhotoURL    string
    IsForwarded bool
}
