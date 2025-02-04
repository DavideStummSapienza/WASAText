package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// ChangeGroupPictureRequest represents the expected JSON body
type ChangeGroupPictureRequest struct {
	NewPhotoURL string `json:"newPhotoURL"`
}

// ChangeGroupPictureResponse represents the JSON response
type ChangeGroupPictureResponse struct {
	Message string `json:"message"`
}

// changeGroupPicture handles updating a group's profile picture
func (rt *_router) changeGroupPicture(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the username from the request context.
	username, ok := r.Context().Value(usernameKey).(string)
	if !ok || username == "" {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Extract group name from the path
	groupName := ps.ByName("groupname")
	if groupName == "" {
		http.Error(w, `{"error": "missing group name in URL path"}`, http.StatusBadRequest)
		return
	}

	// Parse the request body
	var req ChangeGroupPictureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate the new photo URL
	if err := validateURL(req.NewPhotoURL); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Call database function to update group picture
	err := rt.db.ChangeGroupPicture(groupName, req.NewPhotoURL)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Send success response
	response := ChangeGroupPictureResponse{Message: "Group picture successfully updated"}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If there was an error while encoding the response, return a 500 Internal Server Error
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}
