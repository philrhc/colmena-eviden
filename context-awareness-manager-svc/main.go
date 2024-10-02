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
	models "context-awareness-manager/src"
	"context-awareness-manager/src/monitor"
	"context-awareness-manager/src/receiver"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

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

	// Create the roles  table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS dockerRoleDefinitions (
		id TEXT PRIMARY KEY,
		imageId TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/context", func(w http.ResponseWriter, r *http.Request) {
		receiver.HandleContext(w, r, db)
	})

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

	go func() {
		// Listen for changes and print them
		for result := range resultChannel {
			fmt.Printf("Context has changed: %s\n", result.Classification)
			keyExpression := fmt.Sprintf("dockerContextDefinitions/%s", result.ID)
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

// helloHandler provides a welcome message
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Context Awareness Manager Service!")
}

// healthHandler provides a simple health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Context Awareness Manager API is running. Publish new context to /context path")
}
