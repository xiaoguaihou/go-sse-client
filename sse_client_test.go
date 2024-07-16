package client

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Mock server that simulates a server-side event stream.
func setupMockSSEServer() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		// Simulate sending server-side events
		for i := 0; i < 3; i++ {
			event := "data:payload " + strconv.Itoa(i+1) + "\n\n"
			w.Write([]byte(event))
			time.Sleep(time.Second) // Wait for 1 second before sending the next event
		}
	})

	return httptest.NewServer(handler)
}

// Test function for GetSSE.
func TestGetSSE(t *testing.T) {
	mockServer := setupMockSSEServer()
	defer mockServer.Close()

	// Define the test header
	header := map[string]string{
		"Authorization": "Bearer token123",
	}

	// Define the test request
	request := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "John Doe",
		Age:  30,
	}

	// Call GetSSE with the mock server URL, header, request, and mock callback
	err := GetSSE(mockServer.URL, header, request, func(e error, session *map[string]*bytes.Buffer) {
		if e != nil {
			t.Errorf("Error occurred: %v", e)
		} else {
			// Check if the session contains the expected data
			data, ok := (*session)["data:"]
			if ok {
				if strings.Index(data.String(), "payload ") != 0 {
					t.Errorf("Unexpected data: %s", data.String())
				} else {
					fmt.Println(data.String())
				}
			} else {
				t.Errorf("Data not found in session")
			}
		}
	})

	if err != nil {
		t.Errorf("Failed to get SSE: %v", err)
	}
}
