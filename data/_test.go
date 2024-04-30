// main_test.go

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestSuggestionsHandler will test the endpoint /suggestions
func TestSuggestionsHandler(t *testing.T) {
	// Load sample city data for testing
	cities = []City{
		{Name: "Toronto", Latitude: 43.70011, Longitude: -79.4163},
		{Name: "Montreal", Latitude: 45.50884, Longitude: -73.58781},
	}

	// Create a request to pass to the handler
	req, err := http.NewRequest("GET", "/suggestions?q=Toronto", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(suggestionsHandler)

	// Perform the test request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `[{"name":"Toronto","latitude":"43.700110","longitude":"-79.416300","score":1}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// You could also add more tests for different query parameters, error handling, etc.
