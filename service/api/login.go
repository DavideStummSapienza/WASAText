package api

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Define a struct to parse the incoming JSON request body and outgoing response
type LoginRequest struct {
	Username string `json:"username"`
}

type LoginResponse struct {
	Identifier int `json:"identifier"`
}

// doLogin handles user login or account creation.
// If the user exists, their identifier is returned. Otherwise, a new user is created.
func (rt *_router) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Set Content-Type for the response
	w.Header().Set("content-type", "application/json")

	// Parse the JSON request body
	var request LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate the username (length and regex)
	if len(request.Username) < 3 || len(request.Username) > 16 {
		http.Error(w, `{"error": "name must be between 3 and 16 characters"}`, http.StatusBadRequest)
		return
	}

	// Check if the user already exists in the database
	// Check if the username already exists using GetUser
	existingUser, err := rt.db.GetUser(request.Username)
	if err == nil {
		// If user exists, respond with their identifier
		response := LoginResponse{Identifier: existingUser.AuthToken}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	} else if err.Error() != "user not found" {
		// If there's a database error, return an internal server error
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// If the user does not exist, create a new user and return the generated token
	authToken := generateToken() // Generate a random auth token
	rt.db.CreateUser(request.Username, "placeholder", authToken)

	// Respond with the newly created user's auth token
	response := LoginResponse{Identifier: authToken}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// generateToken generates a random integer token
func generateToken() int {
	rand.Seed(time.Now().UnixNano()) // Set the seed for random number generation using the current time
	token := rand.Intn(1000000000)   // Generate a random integer token (between 0 and 999999999)
	return token
}
