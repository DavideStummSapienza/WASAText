package api

import (
	"errors"
	"net/url"
	"strings"
)

// validateURL validates that the given string is a properly formatted URL with http or https schemes.
func validateURL(photoURL string) error {
	// Parse the URL
	parsedURL, err := url.ParseRequestURI(photoURL)
	if err != nil {
		return errors.New("invalid URL format")
	}

	// Check that the scheme is either http or https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("URL must use http or https")
	}

	// Check if the URL likely points to an image
	if !strings.HasSuffix(strings.ToLower(parsedURL.Path), ".jpg") &&
		!strings.HasSuffix(strings.ToLower(parsedURL.Path), ".jpeg") &&
		!strings.HasSuffix(strings.ToLower(parsedURL.Path), ".png") &&
		!strings.HasSuffix(strings.ToLower(parsedURL.Path), ".gif") {
		return errors.New("URL must point to an image file (.jpg, .jpeg, .png, .gif)")
	}

	return nil
}
