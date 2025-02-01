package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DavideStummSapienza/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

// Struct for the decoding the Request
type ForwardMessageRequest struct {
	RecipientUsername string `json:"recipientUsername"`
}

// forwardMessage handles forwarding a message to a new recipient.
//
// Parameters:
// - w: The HTTP response writer used to send responses to the client.
// - r: The HTTP request received from the client.
// - ps: URL parameters extracted by the router.
//
// Behavior:
// - Extracts the authenticated username from the request context.
// - Retrieves the `partner-username` and `message-id` from the URL parameters.
// - Converts `message-id` to an integer and validates it.
// - Decodes the request body to get the new recipient's username.
// - Fetches the original message from the database to ensure it exists and belongs to the user.
// - Forwards the message to the specified new recipient.
// - Responds with appropriate HTTP status codes and messages based on success or failure.
//
// Returns:
// - 200 OK with the forwarded message details if the operation succeeds.
// - 400 Bad Request if required parameters are missing or invalid.
// - 401 Unauthorized if the user is not authenticated.
// - 404 Not Found if the original message does not exist.
// - 500 Internal Server Error if any database operation fails.
func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Set the response content type to JSON.
	w.Header().Set("Content-Type", "application/json")

	// Extract the username from the request context.
	username, ok := r.Context().Value(usernameKey).(string)
	if !ok || username == "" {
		// If the username is missing or invalid, respond with 401 Unauthorized.
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Extract the partner's username from the URL parameters.
	partnerUsername := ps.ByName("partner-username")
	if partnerUsername == "" {
		// If the partnerUsername is missing, respond with 400 Bad Request.
		http.Error(w, `{"error": "partner-username is required"}`, http.StatusBadRequest)
		return
	}

	// Extract the message ID from the URL parameters.
	forwardedMessageIdStr := ps.ByName("message-id")
	if forwardedMessageIdStr == "" {
		// If the message ID is missing, respond with 400 Bad Request.
		http.Error(w, `{"error": "message-id is required"}`, http.StatusBadRequest)
		return
	}

	// Convert message-id from string to integer.
	forwardedMessageId, err := strconv.Atoi(forwardedMessageIdStr)
	if err != nil {
		// If the conversion fails, respond with 400 Bad Request.
		http.Error(w, `{"error": "invalid message-id format"}`, http.StatusBadRequest)
		return
	}

	// Decode the request body to extract the new recipient's username.
	var request ForwardMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		// If JSON decoding fails, respond with 400 Bad Request.
		http.Error(w, `{"error": "invalid request format"}`, http.StatusBadRequest)
		return
	}

	// Retrieve the original message to ensure it exists and belongs to the user.
	originalMessage, err := rt.db.GetMessage(&forwardedMessageId, username, partnerUsername)
	if err != nil {
		// If the original message is not found, respond with 404 Not Found.
		http.Error(w, `{"error": "original message not found"}`, http.StatusNotFound)
		return
	}

	// Create the NewMessage struct to forward the message.
    newMessage := database.NewMessage{
        FromUser:    username,
        ToUser:      request.RecipientUsername,
        Content:     originalMessage.Content,
        IsPhoto:     originalMessage.IsPhoto,
        PhotoURL:    originalMessage.PhotoURL,
        IsForwarded: true, // Mark it as forwarded
    }

	// Forward the message to the new recipient.
	messageID, err := rt.db.SendMessage(newMessage)
	if err != nil {
		// If forwarding fails, respond with 500 Internal Server Error.
		http.Error(w, `{"error": "failed to forward message: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Retrieve the newly forwarded message.
	latestMessage, err := rt.db.GetMessage(&messageID, username, request.RecipientUsername)
	if err != nil {
		// If retrieving the forwarded message fails, respond with 500 Internal Server Error.
		http.Error(w, `{"error": "failed to retrieve forwarded message"}`, http.StatusInternalServerError)
		return
	}

	// Send a 200 OK response with the forwarded message details.
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(latestMessage); err != nil {
		// Handle any potential error during JSON encoding.
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}
