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
	"context-awareness-manager/internal/models"
	"context-awareness-manager/internal/monitor"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// @Summary Receive context
// @Description Endpoint to receive and process Docker context
// @Tags Context
// @Accept  json
// @Produce json
// @Param context body models.Context true "Context to process"
// @Success 200 {object} models.Result "Context processed successfully"
// @Failure 400 {string} string "Invalid context"
// @Router /context [post]
func HandleContext(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var ctx models.ServiceDescription
	// Decode the JSON request body
	err := json.NewDecoder(r.Body).Decode(&ctx)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		log.Printf("Failed to decode JSON: %v", err)
		return
	}

	// Insert new context into the database
	for _, context := range ctx.DockerContextDefinitions {
		_, err = db.Exec(`INSERT OR REPLACE INTO dockerContextDefinitions (id, imageId) VALUES (?, ?)`, context.ID, context.ImageID)
		if err != nil {
			log.Printf("Failed to save context response for ID %s: %v", context.ID, err)
			return
		}
		//LLamar al Container Engine SDK para ejecutar la imagen y que nos devuelva la clasificacion
		classification, err := monitor.DeployContainer(context.ImageID, nil)
		if err != nil {
			log.Printf("Error deploying container: %v\n", err)
			return
		}
		log.Printf("Context classification: %s\n", classification)
	}

	// Insert new role into the database
	for _, role := range ctx.DockerRoleDefinitions {
		_, err = db.Exec(`INSERT OR REPLACE INTO dockerRoleDefinitions (id, imageId) VALUES (?, ?)`, role.ID, role.ImageID)
		if err != nil {
			log.Printf("Failed to save role response for ID %s: %v", role.ID, err)
			return
		}
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Context received: %+v\n", ctx)
}

// HealthHandler checks if the service is up and running.
// @Summary Check API health
// @Description Checks if the service is up and responding.
// @Tags Health
// @Produce text/plain
// @Success 200 {string} string "OK"
// @Router /health [get]
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Context Awareness Manager API is running. Publish new context to /context path")
}
