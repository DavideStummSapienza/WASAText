package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

// Response structure for image upload
type UploadImageResponse struct {
	ImageURL string `json:"imageUrl"`
}

// uploadImage handles the image upload request
func (rt *_router) uploadImage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// Max image size: 10MB
	// Parse the form data, including file
	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		http.Error(w, `{"error": "The uploaded file is too large or invalid."}`, http.StatusBadRequest)
		return
	}

	// Get the file from the form data
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, `{"error": "Unable to retrieve file from form."}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check if the file extension is valid
	allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
	ext := filepath.Ext(handler.Filename)

	if !allowedExtensions[ext] {
		http.Error(w, `{"error": "Invalid file type. Allowed types are: .jpg, .jpeg, .png, .gif."}`, http.StatusBadRequest)
		return
	}

	// Declare the directory where the images will be stored
	uploadDir := "uploads"

	// Check if the directory exists and handle any errors
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {

		// Create the uploads directory if it doesn't exist
		err := os.Mkdir(uploadDir, os.ModePerm)
		if err != nil {
			http.Error(w, `{"error": "Failed to create the uploads directory."}`, http.StatusInternalServerError)
			return
		}
	} else if err != nil {

		// If the error is not because the directory doesn't exist, log the error
		http.Error(w, `{"error": "Failed to check the uploads directory."}`, http.StatusInternalServerError)
		return
	}

	// Create a unique filename for the uploaded image
	newFileName := fmt.Sprintf("%d%s", os.Getpid(), ext)
	filePath := filepath.Join(uploadDir, newFileName)

	// Save the uploaded file on the server
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, `{"error": "Failed to save the file on the server."}`, http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the server
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, `{"error": "Failed to write the file to the server."}`, http.StatusInternalServerError)
		return
	}

	// Return the file path as part of the response
	imageURL := fmt.Sprintf("http://%s/%s", r.Host, filePath)
	response := UploadImageResponse{ImageURL: imageURL}

	// Send the success response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If encoding the response fails, send a 500 Internal Server Error
		http.Error(w, `{"error": "Failed to encode the response."}`, http.StatusInternalServerError)
		return
	}
}
