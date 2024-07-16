# go-sse-client

a product level go sse client withou any other dependencies.
## Usage

refer to sse_client_test.go for example usage.

```go
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
				fmt.Println(data.String())
			} else {
				t.Errorf("Data not found in session")
			}
		}
	})