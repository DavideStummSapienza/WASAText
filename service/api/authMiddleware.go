package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/DavideStummSapienza/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

// Special datatype for context keys to avoid conflicts with other keys in the context.
type contextKey string

const usernameKey contextKey = "username"

// AuthMiddleware is an authentication middleware that validates the Authorization header.
// If the token is valid, it extracts the username and stores it in the request context.
// If the token is missing, invalid, or cannot be verified, it returns an error response.
func AuthMiddleware(db database.AppDatabase, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Check for the presence of the Authorization header.
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			// If the header is missing or improperly formatted, respond with an error.
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error": "missing or invalid token"}`, http.StatusUnauthorized)
			return
		}

		// Extract the token by removing the "Bearer " prefix.
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Look up the username associated with the token in the database.
		username, err := db.GetUsernameByToken(token)
		if err != nil {
			// If there is an error querying the database, respond with an internal server error.
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
			return
		}
		if username == "" {
			// If the token is invalid or does not exist, respond with an unauthorized error.
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// Add the username to the request context for use in downstream handlers.
		ctx := context.WithValue(r.Context(), usernameKey, username)
		r = r.WithContext(ctx)

		// Call the next handler in the chain with the updated request.
		next(w, r, ps)
	}
}
