package api

import (
	"encoding/json"
	"log"
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
// - Retrieves the `partnerUsername` from the route parameters.
// - Marks all messages in the conversation as read for the user.
// - Fetches the conversation details using the database function.
// - Responds with appropriate HTTP status codes and messages for success or failure.
//
// Returns:
// - 200 OK with the conversation details if the operation succeeds.
// - 400 Bad Request if the `partnerUsername` is missing or invalid.
// - 401 Unauthorized if the username is missing or invalid in the context.
// - 404 Not Found if the conversation partner does not exist.
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

	// Extract the partner's username from the URL parameters.
	partnerUsername := ps.ByName("partner-username")
	if partnerUsername == "" {
		// If the partnerUsername is missing, respond with 400 Bad Request.
		http.Error(w, `{"error": "partner-username is required"}`, http.StatusBadRequest)
		return
	}

	log.Printf("INFO: showConversation called for user: %s, partner: %s", username, partnerUsername)

	// Mark all messages in the conversation as received.
	err := rt.db.MarkAllMessagesAsReceived(username, partnerUsername)
	if err != nil {
		log.Printf("ERROR: Failed to mark messages as received: %v", err)
		http.Error(w, `{"error": "failed to update message status"}`, http.StatusInternalServerError)
		return
	}

	// Mark all messages in the conversation as read.
	err = rt.db.MarkAllMessagesAsRead(username, partnerUsername)
	if err != nil {
		log.Printf("ERROR: Failed to mark messages as read: %v", err)
		// If there is an error updating message status, respond with 500 Internal Server Error.
		http.Error(w, `{"error": "failed to update message status "}`, http.StatusInternalServerError)
		return
	}

	// Fetch the conversation details from the database.
	conversation, err := rt.db.ShowConversation(username, partnerUsername)
	if err != nil {
		log.Printf("ERROR: Failed to fetch conversation: %v", err)
		// If there is an error fetching the conversation, respond with 500 Internal Server Error.
		http.Error(w, `{"error": "failed to fetch conversation: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Send a 200 OK response with the conversation details.
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(conversation); err != nil {
		log.Printf("ERROR: Failed to encode response: %v", err)
		// Handle any potential error during JSON encoding.
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}
