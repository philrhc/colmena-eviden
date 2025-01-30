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
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "context-awareness-manager/docs"

	_ "github.com/mattn/go-sqlite3"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Context Awareness Manager API
// @version 1.0
// @description API API for managing Docker context and roles in the COLMENA project.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description "Type 'Bearer TOKEN' to correctly set the API Key"

// @externalDocs.description OpenAPI
// @externalDocs.url https://swagger.io/resources/open-api/

func main() {
	// Create and initialize the database
	db, err := sql.Open("sqlite3", "./context_awareness_manager.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the contexts table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS dockerContextDefinitions (
		id TEXT PRIMARY KEY,
		imageId TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/context", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleContext(w, r, db)
	})
	// Swagger documentation endpoint
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	interval := 30 * time.Second

	// Channel to receive the monitoring results
	resultChannel := make(chan models.Result)

	// Goroutine to check if DockerContextDefinitions table has data and start monitoring
	go func() {
		for {
			contextList, err := monitor.GetContextData(db)
			if err != nil {
				log.Printf("Error checking context data: %v", err)
				time.Sleep(10 * time.Second)
				continue
			}
			if len(contextList) > 0 {
				fmt.Println("Context data found, starting monitoring")
				go monitor.MonitorContext(interval, contextList, resultChannel)
				break
			}
			time.Sleep(10 * time.Second) // Wait before checking again
		}
	}()

	// Get context key expression prefix
	keyExpression_prefix := os.Getenv("CONTEXT_KEY")
	if keyExpression_prefix == "" {
		keyExpression_prefix = "dockerContextDefinitions"
	}

	go func() {
		// Listen for changes and print them
		for result := range resultChannel {
			fmt.Printf("Context has changed: %s\n", result.Classification)
			keyExpression := fmt.Sprintf("%s/%s", keyExpression_prefix, result.ID)
			err := monitor.PublishContext(keyExpression, result.Classification)
			if err != nil {
				fmt.Println("Error publishing context:", err)
			} else {
				fmt.Println("Context published successfully!")
			}
		}
	}()

	fmt.Println("Starting controller on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
