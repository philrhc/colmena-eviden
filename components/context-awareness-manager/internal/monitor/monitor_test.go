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
package monitor

import (
	"context-awareness-manager/internal/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"bou.ke/monkey"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var (
	myDeployContainer = DeployContainer
)

func TestGetContextData(t *testing.T) {
	// Create a new in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	// Create a table for DockerContextDefinitions
	_, err = db.Exec(`CREATE TABLE DockerContextDefinitions (id TEXT PRIMARY KEY, imageId TEXT NOT NULL)`)
	assert.NoError(t, err)

	// Insert some test data
	_, err = db.Exec(`INSERT INTO DockerContextDefinitions (id, imageId) VALUES ('test_id', 'test_image')`)
	assert.NoError(t, err)

	// Call the function
	contextList, err := GetContextData(db)
	assert.NoError(t, err)

	// Verify the result
	assert.Len(t, contextList, 1)
	assert.Equal(t, "test_id", contextList[0].ID)
	assert.Equal(t, "test_image", contextList[0].ImageID)
}

func TestDeployContainer(t *testing.T) {
	// Mock the environment variable for the Docker engine URL
	os.Setenv("DOCKERENGINE_URL", "http://localhost/deploy")

	// Create a mock server to simulate the Docker engine API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req models.DeployRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)

		// Simulate a successful deployment
		response := models.DeployResponse{
			Classification: "classification",
			Error:          "",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Override the Docker engine URL
	os.Setenv("DOCKERENGINE_URL", mockServer.URL)

	// Call the function with mock data
	classification, err := DeployContainer("test_image", nil)
	assert.NoError(t, err)
	assert.Equal(t, "classification", classification)
}

func TestMonitorContext(t *testing.T) {
	// Mock the DeployContainer function using monkey patching
	monkey.Patch(DeployContainer, func(imageID string, cmd []string) (string, error) {
		return "classification", nil
	})
	defer monkey.UnpatchAll() // Ensure to restore the original function after the test

	// Prepare test data
	contextList := []models.DockerContextDefinition{
		{ID: "context_id", ImageID: "test_image"},
	}

	// Create a result channel
	resultChannel := make(chan models.Result, 1)

	// Run MonitorContext in a goroutine
	go MonitorContext(1*time.Second, contextList, resultChannel)

	// Wait for the result or timeout
	select {
	case result := <-resultChannel:
		assert.Equal(t, "context_id", result.ID)
		assert.Equal(t, "classification", result.Classification)
	case <-time.After(2 * time.Second):
		t.Fatal("Did not receive expected result in time")
	}
}

func TestPublishContext(t *testing.T) {
	// Mock the Zenoh URL environment variable
	os.Setenv("ZENOH_URL", "http://localhost/publish/")

	// Create a mock server to simulate the Zenoh HTTP PUT API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/key", r.URL.Path)

		// Simulate a successful PUT request
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	// Override the Zenoh URL
	os.Setenv("ZENOH_URL", mockServer.URL)

	// Call the function
	err := PublishContext("key", "value")
	assert.NoError(t, err)
}
