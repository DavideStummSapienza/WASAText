package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Struct for encoding the response message
type DeleteMessageResponse struct {
	Message string `json:"message"`
}

// deleteMessage handles the deletion of a specific message by ID.
//
// Parameters:
// - w: HTTP response writer
// - r: HTTP request
// - ps: Route parameters (contains message-id)
//
// Returns:
// - 200 OK if the deletion is successful.
// - 400 Bad Request if required parameters are missing/invalid.
// - 401 Unauthorized if the user is not authenticated.
// - 404 Not Found if the message does not exist or the user lacks permissions.
// - 500 Internal Server Error if the deletion process fails.
func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Extract the authenticated username from the request context
	username, ok := r.Context().Value(usernameKey).(string)
	if !ok || username == "" {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// 2. Extract `messageID` from the URL and convert it to an integer
	messageIDStr := ps.ByName("message-id")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		http.Error(w, `{"error": "invalid message-id format"}`, http.StatusBadRequest)
		return
	}

	// 3. Delete the message from the database
	err = rt.db.DeleteMessage(username, messageID)
	if err != nil {
		if err.Error() == "message not found or no permissions to delete" {
			http.Error(w, `{"error": "message not found or not authorized"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "failed to delete message"}`, http.StatusInternalServerError)
		return
	}

	// 4. Return a success response
	response := DeleteMessageResponse{Message: "Message deleted successfully"}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Handle any potential error during JSON encoding.
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}
