package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// listConversations handles requests to retrieve a list of user conversations.
//
// Parameters:
// - w: The HTTP response writer used to send responses to the client.
// - r: The HTTP request received from the client.
// - ps: URL parameters extracted by the router.
//
// Behavior:
// - Extracts the username from the request context (set by the authentication middleware).
// - Loads the user's conversations from the database, including the latest message and metadata.
// - Marks all messages in the conversations as received.
// - Reloads the user's conversations after marking messages as received.
// - Responds with a JSON payload containing the updated list of conversations in reverse chronological order.
//
// Returns:
// - 200 OK and a JSON array of conversations if the operation succeeds.
// - 401 Unauthorized if the username is missing or invalid in the context.
// - 500 Internal Server Error if there is a database error or if the conversations cannot be loaded.
func (rt *_router) listConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extract the username associated with the token from the request context.
	username, ok := r.Context().Value(usernameKey).(string)
	if !ok || username == "" {
		// If the username is missing or invalid, respond with 401 Unauthorized.
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Load the user's conversations from the database to get the partner names.
	conversations, err := rt.db.LoadUserConversations(username)
	if err != nil {
		// If there is an error loading conversations, respond with 500 Internal Server Error.
		http.Error(w, `{"error": "failed to load conversations: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Mark all messages in the retrieved conversations as received.
	for _, conversation := range conversations {
		// For each conversation, update the message status for the current user.
		err := rt.db.MarkAllMessagesAsReceived(conversation.Name, username)
		if err != nil {
			// If there is an error updating message status, respond with 500 Internal Server Error.
			http.Error(w, `{"error": "failed to update message status"}`, http.StatusInternalServerError)
			return
		}
	}

	// Reload the user's conversations after marking messages as received, ensuring we send up-to-date data.
	conversations, err = rt.db.LoadUserConversations(username)
	if err != nil {
		// If there is an error reloading conversations, respond with 500 Internal Server Error.
		http.Error(w, `{"error": "failed to reload conversations: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Set the response content type to JSON.
	w.Header().Set("Content-Type", "application/json")

	// Respond with a 200 OK status and the updated list of conversations as JSON.
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(conversations); err != nil {
		// If there is an error encoding the response, respond with 500 Internal Server Error.
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
	}
}
