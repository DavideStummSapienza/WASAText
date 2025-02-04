package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Response structure for leaving a group
type LeaveGroupResponse struct {
	Message string `json:"message"`
}

// leaveGroup handles the HTTP request to remove a user from a group
func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the username from the request context.
	username, ok := r.Context().Value(usernameKey).(string)
	if !ok || username == "" {
		// If the username is missing or invalid, respond with 401 Unauthorized.
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Get group name from path parameter
	groupName := ps.ByName("groupname")
	if groupName == "" {
		http.Error(w, `{"error": "missing groupname parameter"}`, http.StatusBadRequest)
		return
	}

	// Call the database function to remove the user from the group
	err := rt.db.LeaveGroup(groupName, username)
	if err != nil {
		// If the database function returns an error, respond with a 400 Bad Request error
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Send a success response
	response := LeaveGroupResponse{Message: "Successfully left the group"}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If there was an error while encoding the response, return a 500 Internal Server Error
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}
