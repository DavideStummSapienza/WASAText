package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// ChangeProfilePictureRequest defines the structure of the request payload for changing a profile picture.
type ChangeProfilePictureRequest struct {
	PhotoURL string `json:"photo_url"` // The URL of the new profile picture
}

// ChangeProfilePictureResponse defines the structure of the response payload for a successful profile picture change.
type ChangeProfilePictureResponse struct {
	Message  string `json:"message"`   // Success message
	PhotoURL string `json:"photo_url"` // The updated profile picture URL
}

// changeProfilePicture handles requests to update a user's profile picture.
//
// Parameters:
// - w: The HTTP response writer used to send responses to the client.
// - r: The HTTP request received from the client.
// - ps: URL parameters extracted by the router.
//
// Behavior:
// - Extracts the username from the request context (set by the authentication middleware).
// - Parses the request payload to retrieve the new profile picture URL.
// - Validates the new profile picture URL for correctness.
// - Updates the profile picture URL in the database for the authenticated user.
// - Responds with appropriate HTTP status codes and messages for success or failure.
//
// Returns:
// - 200 OK and a success message with the updated profile picture URL if the operation succeeds.
// - 400 Bad Request if the request body is invalid or the URL is invalid.
// - 401 Unauthorized if the username is missing or invalid in the context.
// - 500 Internal Server Error if there is a database error.
func (rt *_router) changeProfilePicture(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Set the response content type to JSON.
	w.Header().Set("Content-Type", "application/json")

	// Extract the username from the request context.
	username, ok := r.Context().Value(usernameKey).(string)
	if !ok || username == "" {
		// If the username is missing or invalid, respond with 401 Unauthorized.
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Parse the request body into the ChangeProfilePictureRequest struct.
	var request ChangeProfilePictureRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		// If the request body is invalid, respond with 400 Bad Request.
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate the photo URL
	if err := validateURL(request.PhotoURL); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Update the profile picture in the database.
	if err := rt.db.ChangeProfilePicture(username, request.PhotoURL); err != nil {
		// If there is a database error, respond with 500 Internal Server Error.
		http.Error(w, `{"error": "failed to update profile picture: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Create a response with a success message and the updated profile picture URL.
	response := ChangeProfilePictureResponse{
		Message:  "profile picture successfully updated",
		PhotoURL: request.PhotoURL,
	}

	// Send a 200 OK response with the success payload.
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Handle any potential error during JSON encoding.
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}
