package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Request structure for changing a group's name
type ChangeGroupNameRequest struct {
	NewGroupName string `json:"newGroupName"`
}

// Response structure for changing a group's name
type ChangeGroupNameResponse struct {
	Message string `json:"message"`
}

// changeGroupName handles the HTTP request to change a group name
func (rt *_router) changeGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the username from the request context.
	username, ok := r.Context().Value(usernameKey).(string)
	if !ok || username == "" {
		// If the username is missing or invalid, respond with 401 Unauthorized.
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Get old groupname from Path parameter
	oldGroupName := ps.ByName("groupname")
	if oldGroupName == "" {
		http.Error(w, `{"error": "missing groupname parameter"}`, http.StatusBadRequest)
		return
	}

	// Parse the request body
	var req ChangeGroupNameRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// If decoding fails, respond with a 400 Bad Request error
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate the input
	if len(req.NewGroupName) == 0 {
		// If either group name is empty, respond with a 400 Bad Request error
		http.Error(w, `{"error": "new group name is required"}`, http.StatusBadRequest)
		return
	}

	if len(req.NewGroupName) < 3 || len(req.NewGroupName) > 16 {
		http.Error(w, `{"error": "group name must be between 3 and 16 characters"}`, http.StatusBadRequest)
		return
	}

	var ErrUserNotFound = errors.New("user not found")

	_, err = rt.db.GetUser(req.NewGroupName)
	if err == nil {
		http.Error(w, `{"error": "group name already exists as a username"}`, http.StatusBadRequest)
		return
	} else if !errors.Is(err, ErrUserNotFound) {
		http.Error(w, `{"error": "database error while checking username"}`, http.StatusInternalServerError)
		return
	}

	// Call the database function to change the group name
	err = rt.db.ChangeGroupName(oldGroupName, req.NewGroupName)
	if err != nil {
		// If the database function returns an error, respond with a 400 Bad Request error
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Send a success response
	response := ChangeGroupNameResponse{Message: "Group name successfully changed"}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If there was an error while encoding the response, return a 500 Internal Server Error
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}
