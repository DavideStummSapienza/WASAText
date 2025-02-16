package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// MakeCommentResponse represents the response structure after a successful comment creation.
type DeleteCommentResponse struct {
	Message string `json:"message"`
}

// deleteComment handles deleting a user's comment on a message.
//
// Parameters:
// - w: HTTP response writer
// - r: HTTP request
// - ps: Route parameters (contains message ID as a query parameter).
//
// Returns:
// - 200 OK if the comment is successfully deleted.
// - 400 Bad Request if the message ID is missing or invalid.
// - 401 Unauthorized if the user is not authenticated.
// - 500 Internal Server Error if database operation fails.
func (rt *_router) deleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// 1️. Get authenticated username from the context
	username, ok := r.Context().Value(usernameKey).(string)
	if !ok || username == "" {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// 2️. Get message ID from query parameter
	messageIDStr := ps.ByName("message-id")
	if messageIDStr == "" {
		http.Error(w, `{"error": "missing message-id parameter"}`, http.StatusBadRequest)
		return
	}

	// Convert messageID to an integer
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		http.Error(w, `{"error": "invalid message-id format"}`, http.StatusBadRequest)
		return
	}

	// 3️. Attempt to delete the comment from the database
	err = rt.db.DeleteComment(messageID, username)
	if err != nil {
		if err.Error() == "comment not found or user not authorized" {
			http.Error(w, `{"error": "comment not found or unauthorized"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "failed to delete comment"}`, http.StatusInternalServerError)
		return
	}

	// 4. Return success response.
	response := DeleteCommentResponse{Message: "Comment deleted successfully"}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Handle any potential error during JSON encoding.
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}
