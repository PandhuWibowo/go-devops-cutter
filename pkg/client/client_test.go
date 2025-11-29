package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	baseURL := "http://example.com"
	client := NewClient(baseURL)

	if client.BaseURL != baseURL {
		t.Errorf("Expected BaseURL %s, got %s", baseURL, client.BaseURL)
	}

	if client.HTTPClient == nil {
		t.Error("HTTPClient should not be nil")
	}

	if client.HTTPClient.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", client.HTTPClient.Timeout)
	}
}

func TestDoRequestSuccess(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Return success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	body, err := client.doRequest("GET", "/test", nil)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := `{"message":"success"}`
	if string(body) != expected {
		t.Errorf("Expected body %s, got %s", expected, string(body))
	}
}

func TestDoRequestWithBody(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"created"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	requestBody := map[string]string{"name": "test"}
	body, err := client.doRequest("POST", "/create", requestBody)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := `{"status":"created"}`
	if string(body) != expected {
		t.Errorf("Expected body %s, got %s", expected, string(body))
	}
}

func TestDoRequestError(t *testing.T) {
	// Create a test server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"bad request"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.doRequest("GET", "/error", nil)

	if err == nil {
		t.Error("Expected error for 400 status code, got nil")
	}

	expectedError := "API error (400): {\"error\":\"bad request\"}"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDoRequestNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"not found"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.doRequest("GET", "/notfound", nil)

	if err == nil {
		t.Error("Expected error for 404 status code, got nil")
	}
}

func TestDoRequestServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.doRequest("GET", "/error", nil)

	if err == nil {
		t.Error("Expected error for 500 status code, got nil")
	}
}

func TestDoRequestInvalidJSON(t *testing.T) {
	client := NewClient("http://example.com")

	// Test with invalid JSON that can't be marshaled
	invalidBody := make(chan int) // channels can't be marshaled to JSON

	_, err := client.doRequest("POST", "/test", invalidBody)

	if err == nil {
		t.Error("Expected error for invalid JSON body, got nil")
	}
}

func TestDoRequestNetworkError(t *testing.T) {
	// Use an invalid URL to trigger network error
	client := NewClient("http://invalid-domain-that-does-not-exist-12345.com")

	_, err := client.doRequest("GET", "/test", nil)

	if err == nil {
		t.Error("Expected network error, got nil")
	}
}
