package transport

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

// post sends a POST request to the specified host with a JSON body and returns the response body.
// host: URL of the host to send the request to
// data: Byte slice representing the JSON body of the request
// Returns the response body as a byte slice if the request is successful, otherwise an error.
func Post(host string, data []byte) ([]byte, error) {
	// Create a new POST request with the provided host URL and request body
	req, err := http.NewRequest("POST", host, bytes.NewBuffer(data))
	if err != nil {
		// Return an error if request creation fails
		return nil, err
	}

	// Set the content type to JSON, indicating the format of the request body
	req.Header.Set("Content-Type", "application/json")

	// Initialize a new HTTP client to send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Return an error if the request fails
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		// Return an error if the response status code is not OK
		return nil, errors.New("received non-200 response code")
	}

	// Read and return the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		// Return an error if reading the response body fails
		return nil, err
	}

	return responseBody, nil
}
