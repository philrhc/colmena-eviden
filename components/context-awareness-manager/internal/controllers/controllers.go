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
	"context-awareness-manager/internal/monitor"
	"context-awareness-manager/pkg/models"
	"context-awareness-manager/pkg/response"
	"encoding/json"
	"fmt"
	"net/http"
)

type ContextHandler struct {
	monitor monitor.ContextMonitor
}

// NewContextHandler creates a new ContextHandler
func NewContextHandler(monitor monitor.ContextMonitor) *ContextHandler {
	return &ContextHandler{monitor: monitor}
}

// @Summary Receive context
// @Description Endpoint to receive and process Service Description context
// @Tags Context
// @Accept  json
// @Produce json
// @Param context body models.ServiceDescription true "Service with Service Definitions to process"
// @Success 200 {object} models.Result "Context processed successfully"
// @Failure 400 {string} string "Invalid Service Description"
// @Failure 500 {string} string "Internal server error"
// @Router /context [post]
func (c *ContextHandler) HandleContextRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.ERROR(w, http.StatusMethodNotAllowed,
			fmt.Errorf("ERROR invalid request method"))
		return
	}

	// Decode the JSON request body
	var service models.ServiceDescription
	err := json.NewDecoder(r.Body).Decode(&service)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest,
			fmt.Errorf("ERROR reading request body: %v", err))
		return
	}

	// Process the context with the service
	err = c.monitor.RegisterContexts(service.DockerContextDefinitions)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.Success(w, map[string]interface{}{"response": service.DockerContextDefinitions})
}

// HealthHandler checks if the service is up and running.
// @Summary Check API health
// @Description Checks if the service is up and responding.
// @Tags Health
// @Produce text/plain
// @Success 200 {string} string "Context Awareness Manager API is running. Publish new context to /context path"
// @Router /health [get]
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, "Context Awareness Manager API is running. Publish new context to /context path")
}
