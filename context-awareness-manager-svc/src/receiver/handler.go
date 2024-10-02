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
package receiver

import (
	models "context-awareness-manager/src"
	"context-awareness-manager/src/monitor"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// HandleContext handles incoming Context requests
func HandleContext(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var ctx models.Context
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

// HandleSubscription handles incoming Subscription requests
func HandleSubscription(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var sub models.Subscriber

	// Decode the JSON request body
	err := json.NewDecoder(r.Body).Decode(&sub)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		log.Printf("Failed to decode JSON: %v", err)
		return
	}

	// Insert the subscriber into the database
	_, err = db.Exec(`INSERT OR REPLACE INTO subscribers (id, endpoint) VALUES (?, ?)`, sub.ID, sub.Endpoint)
	if err != nil {
		http.Error(w, "Failed to save subscription", http.StatusInternalServerError)
		log.Printf("Failed to save subscription: %v", err)
		return
	}

	// Query for the imageID in DockerContextDefinition
	var imageID string
	err = db.QueryRow(`SELECT imageID FROM DockerContextDefinition WHERE id = ?`, sub.ID).Scan(&imageID)
	if err != nil && err != sql.ErrNoRows {
		// ID not found in the database
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		response := map[string]string{
			"message": fmt.Sprintf("ID %s not found in Database", sub.ID),
		}
		json.NewEncoder(w).Encode(response)
		log.Printf("Context related to ID %s not found!", sub.ID)
		return
	}

	response := models.DockerContextDefinition{
		ID:      sub.ID,
		ImageID: imageID,
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	fmt.Fprintf(w, "Subscription successful: %+v\n", sub)
}

func HandleListSubscriptions(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT id, endpoint FROM subscribers;")
	if err != nil {
		http.Error(w, "Failed to query subscribers", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var subscriptions []models.Subscriber
	for rows.Next() {
		var subscription models.Subscriber
		if err := rows.Scan(&subscription.ID, &subscription.Endpoint); err != nil {
			http.Error(w, "Failed to scan subscription", http.StatusInternalServerError)
			return
		}
		subscriptions = append(subscriptions, subscription)
	}

	response, err := json.Marshal(subscriptions)
	if err != nil {
		http.Error(w, "Failed to marshal subscriptions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
