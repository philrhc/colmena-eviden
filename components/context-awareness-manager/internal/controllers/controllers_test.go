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

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// Mocking the ContextMonitor interface for testing purposes
type MockMonitor struct {
	StoredContexts map[string]models.DockerContextDefinition
}

// Mock implementation of RegisterContexts method
func (m *MockMonitor) RegisterContexts(contexts []models.DockerContextDefinition) error {
	if len(contexts) == 0 {
		return fmt.Errorf("no contexts provided")
	}
	if m.StoredContexts == nil {
		m.StoredContexts = make(map[string]models.DockerContextDefinition)
	}
	for _, ctx := range contexts {
		m.StoredContexts[ctx.ID] = ctx
	}
	return nil
}

func (m *MockMonitor) GetAllContexts() ([]models.DockerContextDefinition, error) {
	contexts := []models.DockerContextDefinition{}
	for _, ctx := range m.StoredContexts {
		contexts = append(contexts, ctx)
	}
	return contexts, nil
}

func (m *MockMonitor) GetContextByID(id string) (*models.DockerContextDefinition, error) {
	ctx, exists := m.StoredContexts[id]
	if !exists {
		return nil, nil
	}
	return &ctx, nil
}

func (m *MockMonitor) DeleteContext(id string) error {
	if _, exists := m.StoredContexts[id]; !exists {
		return fmt.Errorf("context not found")
	}
	delete(m.StoredContexts, id)
	return nil
}

// Mock implementation of StartMonitoring method
func (m *MockMonitor) StartMonitoring(interval time.Duration) {
	// Simulate the monitoring starting...
}

// Test ContextHandler function
func TestContextHandler_Success(t *testing.T) {
	mock := &MockMonitor{StoredContexts: make(map[string]models.DockerContextDefinition)}
	contextHandler := NewContextHandler(mock)

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
	mock := &MockMonitor{StoredContexts: make(map[string]models.DockerContextDefinition)}
	contextHandler := NewContextHandler(mock)

	req, err := http.NewRequest(http.MethodGet, "/context", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	contextHandler.HandleContextRequest(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

// Test bad request (invalid JSON)
func TestContextHandler_InvalidJSON(t *testing.T) {
	mock := &MockMonitor{StoredContexts: make(map[string]models.DockerContextDefinition)}
	contextHandler := NewContextHandler(mock)

	req, err := http.NewRequest(http.MethodPost, "/context", bytes.NewBuffer([]byte("invalid json")))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	contextHandler.HandleContextRequest(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// Test internal server error (when RegisterContexts fails)
func TestContextHandler_InternalError(t *testing.T) {
	mock := &MockMonitor{StoredContexts: make(map[string]models.DockerContextDefinition)}
	contextHandler := NewContextHandler(mock)

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

func TestGetContexts(t *testing.T) {
	mock := &MockMonitor{
		StoredContexts: map[string]models.DockerContextDefinition{
			"ctx-1": {ID: "ctx-1", ImageID: "img-1"},
		},
	}
	handler := NewContextHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/context", nil)
	rr := httptest.NewRecorder()
	handler.GetContexts(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "ctx-1")
}

func TestGetContextByID_Found(t *testing.T) {
	mock := &MockMonitor{
		StoredContexts: map[string]models.DockerContextDefinition{
			"ctx-123": {ID: "ctx-123", ImageID: "img-xyz"},
		},
	}
	handler := NewContextHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/context/ctx-123", nil)
	rr := httptest.NewRecorder()
	handler.GetContextByID(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "img-xyz")
}

func TestGetContextByID_NotFound(t *testing.T) {
	mock := &MockMonitor{StoredContexts: map[string]models.DockerContextDefinition{}}
	handler := NewContextHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/context/unknown-id", nil)
	rr := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"id": "unknown-id"})

	handler.GetContextByID(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeleteContext_Success(t *testing.T) {
	mock := &MockMonitor{
		StoredContexts: map[string]models.DockerContextDefinition{
			"del-1": {ID: "del-1", ImageID: "img-del"},
		},
	}
	handler := NewContextHandler(mock)

	req := httptest.NewRequest(http.MethodDelete, "/context/del-1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "del-1"})
	rr := httptest.NewRecorder()

	handler.DeleteContext(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteContext_NotFound(t *testing.T) {
	mock := &MockMonitor{StoredContexts: map[string]models.DockerContextDefinition{}}
	handler := NewContextHandler(mock)

	req := httptest.NewRequest(http.MethodDelete, "/context/not-found", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "not-found"})
	rr := httptest.NewRecorder()

	handler.DeleteContext(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
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
