package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type SearchRequest struct {
	Username string `json:"username"`
}

// searchUsers handles the search for users based on a partial or full username.
func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Extract the username from the request context.
	currentUser, ok := r.Context().Value(usernameKey).(string)
	if !ok || currentUser == "" {
		// If the username is missing or invalid, respond with 401 Unauthorized.
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Set Content-Type for the response
	w.Header().Set("content-type", "application/json")

	// Parse the query parameters from the URL
	username := r.URL.Query().Get("username") // Get the 'name' query parameter

	// Perform the search using the database
	users, err := rt.db.SearchUser(username, currentUser)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Encode the search result to JSON and send it as response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
	}
}
