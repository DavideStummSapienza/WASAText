package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// showConversation handles requests to fetch the details of a conversation.
//
// Parameters:
// - w: The HTTP response writer used to send responses to the client.
// - r: The HTTP request received from the client.
// - ps: URL parameters extracted by the router.
//
// Behavior:
// - Extracts the username from the request context (set by the authentication middleware).
// - Retrieves the `conversationID` from the route parameters.
// - Fetches the conversation details using the database function.
// - Responds with appropriate HTTP status codes and messages for success or failure.
//
// Returns:
// - 200 OK with the conversation details if the operation succeeds.
// - 400 Bad Request if the `conversationID` is missing or invalid.
// - 401 Unauthorized if the username is missing or invalid in the context.
// - 500 Internal Server Error if there is a database error.
func (rt *_router) showConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Set the response content type to JSON.
	w.Header().Set("Content-Type", "application/json")

	// Extract the username from the request context.
	username, ok := r.Context().Value(usernameKey).(string)
	if !ok || username == "" {
		// If the username is missing or invalid, respond with 401 Unauthorized.
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Extract the conversationID from the URL parameters.
	conversationID := ps.ByName("conversationID")
	if conversationID == "" {
		// If the conversationID is missing, respond with 400 Bad Request.
		http.Error(w, `{"error": "conversationID is required"}`, http.StatusBadRequest)
		return
	}

	// Fetch the conversation details from the database.
	conversation, err := rt.db.ShowConversation(username, conversationID)
	if err != nil {
		// If there is an error fetching the conversation, respond with 500 Internal Server Error.
		http.Error(w, `{"error": "failed to fetch conversation: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Send a 200 OK response with the conversation details.
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(conversation); err != nil {
		// Handle any potential error during JSON encoding.
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}
