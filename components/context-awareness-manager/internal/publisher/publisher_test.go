/*
COLMENA-DESCRIPTION-SERVICE
Copyright © 2024 EVIDEN

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

This work has been implemented within the context of COLMENA project.
*/
package publisher

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublishContextClassification(t *testing.T) {
	// Create a mock server that returns 200 OK
	handler := http.NewServeMux()
	handler.HandleFunc("/testKey", func(w http.ResponseWriter, r *http.Request) {
		// Verify that the method is PUT
		assert.Equal(t, "PUT", r.Method)
		// Read the request body
		body, _ := io.ReadAll(r.Body)
		// Verify the request body
		expectedBody := `{"value":"testValue"}`
		assert.Equal(t, expectedBody, string(body))
		// Return a 200 OK status
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	// Create a zenohPublisher instance with the mock server URL
	publisher := NewPublisher(server.URL)

	// Call PublishContextClassification with test data
	err := publisher.PublishContextClassification("testKey", "testValue")
	assert.NoError(t, err, "PublishContextClassification should not return an error")
}

func TestPublishContextClassification_Error(t *testing.T) {
	// Create a mock server that returns a non-200 status code
	handler := http.NewServeMux()
	handler.HandleFunc("/testKey", func(w http.ResponseWriter, r *http.Request) {
		// Return a 500 Internal Server Error
		w.WriteHeader(http.StatusInternalServerError)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	// Create a zenohPublisher instance with the mock server URL
	publisher := NewPublisher(server.URL)

	// Call PublishContextClassification with test data
	err := publisher.PublishContextClassification("testKey", "testValue")
	assert.Error(t, err, "PublishContextClassification should return an error")
	assert.Contains(t, err.Error(), "request failed with status code", "Error message should contain the status code")
}

func TestDeleteContext(t *testing.T) {
	// Create a mock server that returns 200 OK
	handler := http.NewServeMux()
	handler.HandleFunc("/testKey", func(w http.ResponseWriter, r *http.Request) {
		// Verify that the method is DELETE
		assert.Equal(t, "DELETE", r.Method)
		// Return a 200 OK status
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	// Create a zenohPublisher instance with the mock server URL
	publisher := NewPublisher(server.URL)

	// Call DeleteContext with test data
	err := publisher.DeleteContext("testKey")
	assert.NoError(t, err, "DeleteContext should not return an error")
}

func TestDeleteContext_Error(t *testing.T) {
	// Create a mock server that returns a non-200 status code
	handler := http.NewServeMux()
	handler.HandleFunc("/testKey", func(w http.ResponseWriter, r *http.Request) {
		// Return a 500 Internal Server Error
		w.WriteHeader(http.StatusInternalServerError)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	// Create a zenohPublisher instance with the mock server URL
	publisher := NewPublisher(server.URL)

	// Call DeleteContext with test data
	err := publisher.DeleteContext("testKey")
	assert.Error(t, err, "DeleteContext should return an error")
	assert.Contains(t, err.Error(), "request failed with status code", "Error message should contain the status code")
}
