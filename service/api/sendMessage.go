package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/DavideStummSapienza/WASAText/service/database"
)


// sendMessageRequest defines the structure of the request payload for sending a message.
type sendMessageRequest struct {
	Message   string `json:"message"`    // The message content (text)
	IsPhoto   bool   `json:"is_photo"`    // Whether the message is a photo message
	PhotoURL  string `json:"photo_url"`   // URL of the photo (if the message is a photo)
}


// sendMessage handles the API request to send a message to a specific user or group.
//
// Parameters:
// - w: The HTTP response writer used to send responses to the client.
// - r: The HTTP request received from the client.
// - ps: URL parameters extracted by the router.
//
// Behavior:
// - Extracts the username from the request context (set by the authentication middleware).
// - Retrieves the `partnerUsername` from the URL parameters.
// - Sends a message to the specified partner and creates a new conversation if necessary.
// - Responds with the details of the sent message or an error.
//
// Returns:
// - 200 OK and the message details if the operation succeeds.
// - 400 Bad Request if the `message` or `partner-username` is missing or invalid.
// - 401 Unauthorized if the username is missing or invalid in the context.
// - 404 Not Found if the conversation partner does not exist or is not a valid user.
// - 500 Internal Server Error if there is a database error or if the message cannot be sent.
 
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Set the content type of the response to JSON.
	w.Header().Set("Content-Type", "application/json")

	// Extract the username from the request context (set by authentication middleware).
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

	// Parse the request body to retrieve the message details.
	var request sendMessageRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		// If there's an error decoding the request body, respond with 400 Bad Request.
		http.Error(w, `{"error": "invalid message format"}`, http.StatusBadRequest)
		return
	}

	// Validate that the message or photoUrl is provided.
	if request.Message == "" && !request.IsPhoto {
		http.Error(w, `{"error": "message or photoUrl is required"}`, http.StatusBadRequest)
		return
	}

	// Validate PhotoURL if IsPhoto is true
	if request.IsPhoto {
		if err := validateURL(request.PhotoURL); err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}
	}

	// Create the message struct
	newMessage := database.NewMessage{
		FromUser:    username,
		ToUser:      partnerUsername,
		Content:     request.Message,
		IsPhoto:     request.IsPhoto,
		PhotoURL:    request.PhotoURL,
		IsForwarded: false, // Standardmäßig nicht weitergeleitet
	}

	// Create the message in the database by calling SendMessage
	messageID, err := rt.db.SendMessage(newMessage)
	if err != nil {
		// If there is an error sending the message, respond with 500 Internal Server Error.
		http.Error(w, `{"error": "failed to send message: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Retrieve the conversation details including the message.
	latestMessage, err := rt.db.GetMessage(&messageID,username,partnerUsername)
	if err != nil {
		// If there's an error fetching the conversation details, respond with 500 Internal Server Error.
		http.Error(w, `{"error": "failed to retrieve conversation detail: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Respond with the sent message details.
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(latestMessage); err != nil {
		// If there's an error encoding the response, respond with 500 Internal Server Error.
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}

