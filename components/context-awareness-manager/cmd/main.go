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
	"context-awareness-manager/pkg/server"
	"os"
	"sync"

	_ "context-awareness-manager/docs"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
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

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description "Type 'Bearer TOKEN' to correctly set the API Key"

// @externalDocs.description OpenAPI
// @externalDocs.url https://swagger.io/resources/open-api/

func main() {
	// Get environment variables with defaults
	zenohURL := getEnv("ZENOH_URL", "http://zenoh-router:8000")
	port := getEnv("SERVER_PORT", "8080")

	// Initialize router
	router, contextMonitor, dbConnection := server.InitRouter(zenohURL)

	// WaitGroup to handle multiple goroutines
	var wg sync.WaitGroup

	// Start the server in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done() // Decrement the WaitGroup counter once done
		logrus.Infof("Starting server on %s...", port)
		server.StartServer(port, router) // Start the server
	}()

	// Start context monitoring in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done() // Decrement the WaitGroup counter once done
		logrus.Infof("Starting context monitoring...")
		contextMonitor.StartMonitoring()
	}()

	// Wait for all goroutines to finish before exiting the application
	wg.Wait()
	defer dbConnection.Close() // Ensure the database connection is closed at the end
}

// getEnv is a helper function to retrieve environment variables with defaults
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
