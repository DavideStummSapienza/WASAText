package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// MakeCommentRequest represents the expected structure of the request body.
type MakeCommentRequest struct {
	Content string `json:"content"`
}

// MakeCommentResponse represents the response structure after a successful comment creation.
type MakeCommentResponse struct {
	Message string `json:"message"`
}

// makeComment handles adding a comment to a specific message.
//
// Parameters:
// - w: HTTP response writer.
// - r: HTTP request.
// - ps: Route parameters containing the message ID.
//
// Responses:
// - 201 Created: The comment was successfully added.
// - 400 Bad Request: The request is missing required parameters or has an invalid format.
// - 401 Unauthorized: The user is not authenticated.
// - 404 Not Found: The referenced message does not exist.
// - 500 Internal Server Error: A database error occurred.
func (rt *_router) makeComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// Retrieve the authenticated user from the request context.
	username, ok := r.Context().Value(usernameKey).(string)
	if !ok || username == "" {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Extract and validate the message ID from the URL parameters.
	messageIDStr := ps.ByName("message-id")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		http.Error(w, `{"error": "invalid message-id format"}`, http.StatusBadRequest)
		return
	}

	// Decode and validate the request body.
	var req MakeCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		http.Error(w, `{"error": "invalid or missing comment content"}`, http.StatusBadRequest)
		return
	}

	// Insert the comment into the database.
	err = rt.db.AddComment(messageID, username, req.Content)
	if err != nil {
		http.Error(w, `{"error": "failed to add comment"}`, http.StatusInternalServerError)
		return
	}

	// Return success response.
	response := MakeCommentResponse{Message: "Comment added successfully"}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Handle any potential error during JSON encoding.
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}
