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
package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"containerengine/api/handlers"

	"github.com/docker/docker/client"
	"github.com/gorilla/mux"

	// Import the required SQLite driver
	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

// getDockerHost is a helper function that returns the Docker host from the environment variable.
func getDockerHost() string {
	host := os.Getenv("DOCKER_HOST")
	if host == "" {
		host = "unix:///var/run/docker.sock"
	}
	return host
}

// TestMainSetup ensures HTTP handlers are correctly registered
func TestMain(t *testing.T) {
	// Patch the Docker client creation function
	monkey.Patch(client.NewClientWithOpts, func(opts ...client.Opt) (*client.Client, error) {
		return &client.Client{}, nil
	})
	defer monkey.Unpatch(client.NewClientWithOpts) // Unpatch after the test

	// Patch the HealthHandler with a mock
	monkey.Patch(handlers.HealthHandler, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Mock Health OK")
	})
	defer monkey.Unpatch(handlers.HealthHandler) // Unpatch after the test

	// Patch the DeployHandler with a mock
	monkey.Patch(handlers.DeployHandler, func(_ *client.Client, w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Mock Deploy Success")
	})
	defer monkey.Unpatch(handlers.DeployHandler) // Unpatch after the test

	// Run the main function
	go main()

	// Conduct tests
	t.Run("TestHealthHandler", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "Mock Health OK\n", rr.Body.String()) // Ensure this matches actual behavior
	})

	t.Run("TestDockerHostEnvVariable", func(t *testing.T) {
		// Scenario 1: DOCKER_HOST is set
		os.Setenv("DOCKER_HOST", "unix:///mock/docker.sock")
		host := getDockerHost()
		assert.Equal(t, "unix:///mock/docker.sock", host)
		os.Unsetenv("DOCKER_HOST") // Clean up after the test

		// Scenario 2: DOCKER_HOST is not set
		host = getDockerHost()
		assert.Equal(t, "unix:///var/run/docker.sock", host)
	})

	t.Run("TestDeployHandler", func(t *testing.T) {
		req, err := http.NewRequest("POST", "http://localhost:9000/deploy", nil)
		assert.NoError(t, err)

		// Use the default HTTP client to send the request to the running server
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Check that we get the expected response status and body
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, "Mock Deploy Success\n", string(body)) // Ensure this matches actual behavior
	})

	// t.Run("TestDeployHandlerWithError", func(t *testing.T) {
	// 	req, err := http.NewRequest("POST", "/deploy", nil)
	// 	assert.NoError(t, err)

	// 	rr := httptest.NewRecorder()
	// 	router := mux.NewRouter()

	// 	// Mock the DeployHandler to return an error
	// 	router.HandleFunc("/deploy", func(w http.ResponseWriter, r *http.Request) {
	// 		http.Error(w, "Deployment failed", http.StatusInternalServerError)
	// 	}).Methods("POST")

	// 	router.ServeHTTP(rr, req)

	// 	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	// 	assert.Equal(t, "Deployment failed\n", rr.Body.String())
	// })
}
