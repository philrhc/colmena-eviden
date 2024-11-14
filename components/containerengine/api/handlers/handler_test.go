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
package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"containerengine/internal/dockerclient"
	"containerengine/internal/models"

	"bou.ke/monkey"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Container Engine API is running. Publish new deployment to /deploy path\n"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestDeployHandler tests the DeployHandler using monkey patching
func TestDeployHandler(t *testing.T) {
	// Test case: Valid request
	t.Run("Valid request", func(t *testing.T) {
		// Prepare a mock of RunContainer using monkey patching
		monkey.Patch(dockerclient.RunContainer, func(cli *client.Client, image string, cmd []string) (string, error) {
			return "container logs", nil // Simulate a successful response
		})
		defer monkey.UnpatchAll() // Make sure to unpatch all changes at the end

		// Create the router and request
		router := http.NewServeMux()
		router.HandleFunc("/deploy", func(w http.ResponseWriter, r *http.Request) {
			DeployHandler(nil, w, r) // You can pass nil since we're mocking
		})

		reqBody := models.DeployRequest{
			Image: "test-image",
			Cmd:   []string{"arg1", "arg2"},
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/deploy", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response models.DeployResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "container logs", response.Classification)
	})

	// Test case: Invalid request
	t.Run("Invalid request payload", func(t *testing.T) {
		router := http.NewServeMux()
		router.HandleFunc("/deploy", func(w http.ResponseWriter, r *http.Request) {
			DeployHandler(nil, w, r) // You can pass nil since we're mocking
		})

		req := httptest.NewRequest(http.MethodPost, "/deploy", bytes.NewBufferString("invalid json"))
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// Test case: Error while running the container
	t.Run("RunContainer error", func(t *testing.T) {
		// Prepare a mock of RunContainer that returns an error
		monkey.Patch(dockerclient.RunContainer, func(cli *client.Client, image string, cmd []string) (string, error) {
			return "", fmt.Errorf("some error message")
		})
		defer monkey.UnpatchAll() // Make sure to unpatch all changes at the end

		router := http.NewServeMux()
		router.HandleFunc("/deploy", func(w http.ResponseWriter, r *http.Request) {
			DeployHandler(nil, w, r) // You can pass nil since we're mocking
		})

		reqBody := models.DeployRequest{
			Image: "test-image",
			Cmd:   []string{"arg1", "arg2"},
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/deploy", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var response models.DeployResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "some error message", response.Error)
	})
}
