package transport

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestPost includes subtests to cover various scenarios in a single test function
func TestPost(t *testing.T) {
	// Subtest for a successful 200 OK response
	t.Run("Success", func(t *testing.T) {
		// Set up a test server that responds with 200 OK and a JSON body
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "success"}`))
		}))
		defer testServer.Close()

		// Call the post function with the test server URL
		resp, err := Post(testServer.URL, []byte(`{"data": "test"}`))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify the response body matches the expected output
		expected := `{"status": "success"}`
		if string(resp) != expected {
			t.Errorf("Expected %s, got %s", expected, resp)
		}
	})

	// Subtest for a non-200 response code
	t.Run("NonOKResponse", func(t *testing.T) {
		// Set up a test server that responds with 500 Internal Server Error
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer testServer.Close()

		// Call the post function with the test server URL
		_, err := Post(testServer.URL, []byte(`{"data": "test"}`))
		if err == nil {
			t.Fatal("Expected error for non-200 response code, got nil")
		}
	})

	// Subtest for an error reading the response body
	t.Run("ReadBodyError", func(t *testing.T) {
		// Set up a test server that closes the connection immediately
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "1") // incorrect content length to simulate read error
		}))
		defer testServer.Close()

		// Call the post function with the test server URL
		_, err := Post(testServer.URL, []byte(`{"data": "test"}`))
		if err == nil {
			t.Fatal("Expected error due to body read failure, got nil")
		}
	})

	// Subtest for an error when creating the request
	t.Run("RequestCreationError", func(t *testing.T) {
		// Provide an invalid URL to simulate a request creation error
		_, err := Post(":", []byte(`{"data": "test"}`)) // invalid URL
		if err == nil {
			t.Fatal("Expected error due to request creation failure, got nil")
		}
	})
}
