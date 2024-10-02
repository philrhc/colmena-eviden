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
	"bytes"
	models "context-awareness-manager/src"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// GetContextData check if DockerContextDefinitions table has data and return imageID list
func GetContextData(db *sql.DB) ([]models.DockerContextDefinition, error) {
	rows, err := db.Query("SELECT id, imageId FROM DockerContextDefinitions")
	if err != nil {
		return nil, fmt.Errorf("failed to query DockerContextDefinitions table: %v", err)
	}
	defer rows.Close()

	var contextList []models.DockerContextDefinition
	for rows.Next() {
		var context models.DockerContextDefinition
		if err := rows.Scan(&context.ID, &context.ImageID); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		contextList = append(contextList, context)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}
	return contextList, nil
}

func DeployContainer(image string, cmd []string) (string, error) {
	// Retrieve the destination URL from an environment variable
	destinationURL := os.Getenv("DOCKERENGINE_URL")
	if destinationURL == "" {
		log.Fatal("DOCKERENGINE_URL environment variable is not set")
	}
	// Build the request
	requestBody := models.DeployRequest{
		Image: image,
		Cmd:   cmd,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON request: %v", err)
	}

	// Make the HTTP POST request to the microservice
	resp, err := http.Post(destinationURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error making HTTP POST request: %v", err)
	}
	defer resp.Body.Close()

	// Process the response
	var response models.DeployResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", fmt.Errorf("error decoding JSON response: %v", err)
	}

	if response.Error != "" {
		return "", fmt.Errorf("deployment error: %s", response.Error)
	}

	return response.Classification, nil
}

// MonitorContext makes the HTTP POST request and compare the results
func MonitorContext(interval time.Duration, contextList []models.DockerContextDefinition, resultChannel chan models.Result) {
	var lastResults = make(map[string]string)

	for {
		for _, context := range contextList {
			currentClassification, err := DeployContainer(context.ImageID, nil)
			if err != nil {
				log.Printf("Error deploying container: %v\n", err)
				continue
			}
			log.Printf("Context classification: %s\n", currentClassification)

			if lastResult, ok := lastResults[context.ImageID]; !ok || currentClassification != lastResult {
				resultChannel <- models.Result{ID: context.ID, Classification: currentClassification}
				lastResults[context.ImageID] = currentClassification
			}

			time.Sleep(interval)
		}
	}
}

// PublishContext sends an HTTP PUT request to the specified URL with a JSON body
// containing the provided value. The function takes two string parameters: `url` and `value`.
// It returns an error if the PUT request fails or if the response status code is not 200 OK.
func PublishContext(key string, value string) error {
	// Construct the URL using the provided key
	zenohURL := os.Getenv("ZENOH_URL")
	if zenohURL == "" {
		zenohURL = "http://zenoh:8000/"
	}
	url := fmt.Sprintf("%s%s", zenohURL, key)

	// Create the JSON body with the provided value
	data := map[string]string{"value": value}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error converting data to JSON: %v", err)
	}

	log.Printf("Send %s to %s", jsonData, url)

	// Create a new HTTP PUT request
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating PUT request: %v", err)
	}

	// Set the Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP PUT request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending PUT request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	return nil
}
