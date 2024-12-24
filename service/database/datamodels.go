package database

import "time"

// ConversationPreview represents the preview information for a conversation.
type ConversationPreview struct {
	Name            string    `json:"name"`              // Username or Group Name
	PhotoURL        string    `json:"photo_url"`         // Profile Photo URL
	LastMessage     string    `json:"last_message"`      // Snippet or Photo Icon
	LastMessageTime time.Time `json:"last_message_time"` // Timestamp of Last Message
}
