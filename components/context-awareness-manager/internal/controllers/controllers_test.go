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
package controllers

import (
	"bytes"
	"context-awareness-manager/pkg/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Mocking the ContextMonitor interface for testing purposes
type MockMonitor struct{}

// Mock implementation of RegisterContexts method
func (m *MockMonitor) RegisterContexts(contexts []models.DockerContextDefinition) error {
	if len(contexts) == 0 {
		return fmt.Errorf("no contexts provided")
	}
	return nil
}

// Mock implementation of StartMonitoring method
func (m *MockMonitor) StartMonitoring(interval time.Duration) {
	// Simulate the monitoring starting...
}

// Test ContextHandler function
func TestContextHandler_Success(t *testing.T) {
	mockMonitor := &MockMonitor{}
	contextHandler := NewContextHandler(mockMonitor)

	service := models.ServiceDescription{
		DockerContextDefinitions: []models.DockerContextDefinition{{ID: "test-id", ImageID: "test-image"}},
	}
	body, _ := json.Marshal(service)

	req := httptest.NewRequest(http.MethodPost, "/context", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	// Call the handler
	contextHandler.HandleContextRequest(rr, req)

	resp := rr.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// Test invalid HTTP method
func TestContextHandler_InvalidMethod(t *testing.T) {
	mockMonitor := &MockMonitor{}
	contextHandler := NewContextHandler(mockMonitor)

	req, err := http.NewRequest(http.MethodGet, "/context", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	contextHandler.HandleContextRequest(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

// Test bad request (invalid JSON)
func TestContextHandler_InvalidJSON(t *testing.T) {
	mockMonitor := &MockMonitor{}
	contextHandler := NewContextHandler(mockMonitor)

	req, err := http.NewRequest(http.MethodPost, "/context", bytes.NewBuffer([]byte("invalid json")))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	contextHandler.HandleContextRequest(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// Test internal server error (when RegisterContexts fails)
func TestContextHandler_InternalError(t *testing.T) {
	mockMonitor := &MockMonitor{}
	contextHandler := NewContextHandler(mockMonitor)

	service := models.ServiceDescription{
		DockerContextDefinitions: []models.DockerContextDefinition{}, // no contexts, should trigger an error
	}
	body, _ := json.Marshal(service)

	req := httptest.NewRequest(http.MethodPost, "/context", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	// Call the handler
	contextHandler.HandleContextRequest(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "no contexts provided")
}

// Test HealthHandler
func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	HealthHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	expectedResponse := `"Context Awareness Manager API is running. Publish new context to /context path"`
	assert.JSONEq(t, expectedResponse, rr.Body.String())
}
