package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// ChangeUsernameRequest defines the structure of the request payload for changing a username.
type ChangeUsernameRequest struct {
	NewUsername string `json:"newusername"` // The desired new username
}

// ChangeUsernameResponse defines the structure of the response payload for a successful username change.
type ChangeUsernameResponse struct {
	Message     string `json:"message"`     // Success message
	NewUsername string `json:"newusername"` // The updated username
}

// changeUsername handles requests to change a user's username.
//
// Parameters:
// - w: The HTTP response writer used to send responses to the client.
// - r: The HTTP request received from the client.
// - ps: URL parameters extracted by the router.
//
// Behavior:
// - Validates the request payload for correctness.
// - Checks if the new username meets the required criteria.
// - Ensures the new username does not already exist in the database.
// - Updates the username in the database if all checks pass.
// - Responds with appropriate HTTP status codes and messages for success or failure.
//
// Returns:
// - 200 OK and a success message with the updated username if the operation succeeds.
// - 400 Bad Request if the request body is invalid, the username is invalid, or the username already exists.
// - 500 Internal Server Error if there is an unexpected database error.
func (rt *_router) changeUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Set the content type of the response to JSON
	w.Header().Set("content-type", "application/json")

	// Extract the username from the request context.
	oldUsername, ok := r.Context().Value(usernameKey).(string)
	if !ok || oldUsername == "" {
		// If the username is missing or invalid, respond with 401 Unauthorized.
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var request ChangeUsernameRequest

	// Parse the JSON request body into the ChangeUsernameRequest struct
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate the new username: must be between 3 and 16 characters
	if len(request.NewUsername) < 3 || len(request.NewUsername) > 16 {
		http.Error(w, `{"error": "name must be between 3 and 16 characters"}`, http.StatusBadRequest)
		return
	}

	var ErrUserNotFound = errors.New("user not found")
	var ErrGroupNotFound = errors.New("group not found")

	// Check if the new username already exists in the database
	_, err := rt.db.GetUser(request.NewUsername)
	if err == nil {
		// If no error, the username already exists
		http.Error(w, `{"error": "username already exists"}`, http.StatusBadRequest)
		return
	} else if !errors.Is(err,ErrUserNotFound) {
		// If an unexpected database error occurs
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	_, err = rt.db.GetGroupByName(request.NewUsername)
	if err == nil {
		http.Error(w, `{"error": "username cannot be the same as a group name"}`, http.StatusBadRequest)
		return
	} else if !errors.Is(err, ErrGroupNotFound) {
		http.Error(w, `{"error": "database error"}`, http.StatusInternalServerError)
		return
	}

	// Attempt to update the username in the database
	if err := rt.db.ChangeUsername(oldUsername, request.NewUsername); err != nil {
		// If there is a database error while changing the username
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Create a response with a success message and the updated username
	response := ChangeUsernameResponse{
		Message:     "username successfully changed",
		NewUsername: request.NewUsername,
	}

	// Send a 200 OK response with the success payload
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Handle any potential error during JSON encoding.
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}
