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
	"context-awareness-manager/internal/models"
	"context-awareness-manager/internal/monitor"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bou.ke/monkey"
	_ "github.com/mattn/go-sqlite3"
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
	expected := "Context Awareness Manager API is running. Publish new context to /context path\n"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestHandleContext tests the DeployHandler using monkey patching
func TestHandleContext(t *testing.T) {
	// Test case: Valid request
	t.Run("Valid request", func(t *testing.T) {
		// Create a new in-memory SQLite database for testing
		db, err := sql.Open("sqlite3", ":memory:")
		if err != nil {
			t.Fatalf("Could not open database: %v", err)
		}
		defer db.Close()

		// Create the tables
		_, err = db.Exec(`CREATE TABLE dockerContextDefinitions (id TEXT PRIMARY KEY, imageId TEXT NOT NULL)`)
		if err != nil {
			t.Fatalf("Could not create table: %v", err)
		}
		_, err = db.Exec(`CREATE TABLE dockerRoleDefinitions (id TEXT PRIMARY KEY, imageId TEXT NOT NULL)`)
		if err != nil {
			t.Fatalf("Could not create table: %v", err)
		}

		// Define a valid context based on the provided JSON
		ctx := models.ServiceDescription{
			ID: models.ID{Value: "ExampleApplication"},
			DockerContextDefinitions: []models.DockerContextDefinition{
				{ID: "company_premises", ImageID: "xaviercasasbsc/company_premises"},
			},
			KPIs:                  []models.KPI{},
			DockerRoleDefinitions: []models.DockerRoleDefinition{},
		}

		// Marshal the context to JSON
		jsonData, err := json.Marshal(ctx)
		if err != nil {
			t.Fatalf("Could not marshal context: %v", err)
		}

		req, err := http.NewRequest("POST", "/context", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Could not create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		// Mock the DeployContainer function using monkey patching
		monkey.Patch(monitor.DeployContainer, func(imageID string, cmd []string) (string, error) {
			return "classification", nil
		})
		defer monkey.UnpatchAll() // Ensure to restore the original function after the test

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			HandleContext(w, r, db)
		})

		handler.ServeHTTP(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusOK, rr.Code)

		// Check the response body
		expected := `Context received: {ID:{Value:ExampleApplication} DockerContextDefinitions:[{ID:company_premises ImageID:xaviercasasbsc/company_premises}] KPIs:[] DockerRoleDefinitions:[]}` + "\n"
		assert.Equal(t, expected, rr.Body.String())
	})
}
