package api

import (
	"net/http"
	"strings"
)


// AuthMiddleware checks Authentificationtoken
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            w.Header().Set("Content-Type", "application/json")
            http.Error(w, `{"error": "missing or invalid token"}`, http.StatusUnauthorized)
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        if !isValidToken(token) {
            w.Header().Set("Content-Type", "application/json")
            http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}


// Dummy-Function to validate Token
func isValidToken(token string) bool {
	// Implement real check!!!!!!!!!!!!!!!!!!!!!!!!!
	return token == "expected_token"
}
