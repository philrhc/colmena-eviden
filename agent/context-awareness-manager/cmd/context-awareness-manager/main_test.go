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
	"context-awareness-manager/api/handlers"
	"context-awareness-manager/internal/models"
	"context-awareness-manager/internal/monitor"
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	// Import the required SQLite driver
	"bou.ke/monkey"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

// setupDatabase initializes the in-memory database and creates required tables.
func setupDatabase() (*sql.DB, error) {
	// Create an in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	// Create the contexts table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS dockerContextDefinitions (
		id TEXT PRIMARY KEY,
		imageId TEXT NOT NULL
	)`)
	if err != nil {
		return nil, err
	}

	// Create the roles table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS dockerRoleDefinitions (
		id TEXT PRIMARY KEY,
		imageId TEXT NOT NULL
	)`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// TestDatabaseInitialization ensures the database is set up correctly
func TestDatabaseInitialization(t *testing.T) {
	// Initialize the database
	db, err := setupDatabase()
	if err != nil {
		t.Fatalf("Could not initialize database: %v", err)
	}
	defer db.Close()

	// Check the database connection
	err = db.Ping()
	if err != nil {
		t.Fatalf("Could not ping database: %v", err)
	}

	// Verify that the tables exist
	var count int
	err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='dockerContextDefinitions'").Scan(&count)
	if err != nil || count == 0 {
		t.Fatalf("dockerContextDefinitions table does not exist: %v", err)
	}

	err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='dockerRoleDefinitions'").Scan(&count)
	if err != nil || count == 0 {
		t.Fatalf("dockerRoleDefinitions table does not exist: %v", err)
	}
}

// TestHttpHandlers ensures HTTP handlers are correctly registered
func TestHttpHandlers(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure that the handlers work
		switch r.URL.Path {
		case "/health":
			handlers.HealthHandler(w, r)
		case "/context":
			// Mock context handling (you can expand this further)
			handlers.HandleContext(w, r, nil) // Pass nil as db for simplicity in this test
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Test the /health endpoint
	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Could not get health endpoint: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}
}

// TestMonitoringGoroutine tests the goroutine for monitoring Docker context.
func TestMonitoringGoroutine(t *testing.T) {
	// Initialize the database
	db, err := setupDatabase()
	if err != nil {
		t.Fatalf("Could not initialize database: %v", err)
	}
	defer db.Close()

	// Mocking GetContextData to return a valid context
	monkey.Patch(monitor.GetContextData, func(db *sql.DB) ([]models.DockerContextDefinition, error) {
		return []models.DockerContextDefinition{
			{ID: "context", ImageID: "test-image-id"},
		}, nil // No error
	})
	defer monkey.Unpatch(monitor.GetContextData) // Restore original function after the test

	// Mocking MonitorContext to send results to the channel
	resultChannel := make(chan models.Result)
	monkey.Patch(monitor.MonitorContext, func(interval time.Duration, contexts []models.DockerContextDefinition, resultCh chan models.Result) {
		// Simulate sending a result after monitoring
		for _, ctx := range contexts {
			// Simulate processing and then send a result
			time.Sleep(50 * time.Millisecond)
			resultCh <- models.Result{ID: ctx.ID, Classification: "classification"}
		}
		// Close the channel after sending results
		close(resultCh)
	})
	defer monkey.Unpatch(monitor.MonitorContext) // Restore original function after the test

	var wg sync.WaitGroup
	wg.Add(1)

	// Start the monitoring goroutine
	interval := 30 * time.Second
	go func() {
		defer wg.Done() // Mark this goroutine as done when finished
		for {
			contextList, err := monitor.GetContextData(db)
			if err != nil {
				log.Printf("Error checking context data: %v", err)
				time.Sleep(10 * time.Second)
				continue
			}
			if len(contextList) > 0 {
				log.Println("Context data found, starting monitoring")
				go monitor.MonitorContext(interval, contextList, resultChannel)
				break
			}
			time.Sleep(10 * time.Second) // Wait before checking again
		}
	}()

	// Check if the resultChannel received a result
	select {
	case result := <-resultChannel:
		assert.Equal(t, "classification", result.Classification, "Expected classification to match")
	case <-time.After(500 * time.Millisecond): // Timeout to prevent hanging
		t.Error("Expected result channel to receive a result, but it was empty")
	}

	// Wait for the goroutine to finish
	wg.Wait()
}

// TestPublishingGoroutine tests the goroutine for publishing Docker context.
func TestPublishingGoroutine(t *testing.T) {
	// Create a buffered channel to receive monitoring results
	resultChannel := make(chan models.Result, 2)

	// Mock PublishContext to track calls
	mockPublishCalled := false
	// Mocking PublishContext to do nothing
	monkey.Patch(monitor.PublishContext, func(key string, classification string) error {
		mockPublishCalled = true // Mark that it was called
		return nil               // Simulate successful publishing
	})
	defer monkey.Unpatch(monitor.PublishContext) // Restore original function after the test

	var wg sync.WaitGroup
	wg.Add(1)

	// Goroutine to listen for results
	go func() {
		defer wg.Done() // Mark this goroutine as done when finished
		for result := range resultChannel {
			log.Printf("Context has changed: %s\n", result.Classification)
			keyExpression := "dockerContextDefinitions/" + result.ID
			err := monitor.PublishContext(keyExpression, result.Classification)
			assert.NoError(t, err, "Expected no error while publishing context")
		}
	}()

	// Simulate sending results to the resultChannel
	resultChannel <- models.Result{ID: "context-1", Classification: "classification"}
	close(resultChannel) // Close channel after sending results
	// Wait for the goroutine to finish
	wg.Wait()

	// Assert that PublishContext was called
	assert.True(t, mockPublishCalled, "Expected PublishContext to be called")
}
