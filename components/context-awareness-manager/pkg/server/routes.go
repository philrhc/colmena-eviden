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
package server

import (
	"context-awareness-manager/internal/controllers"
	"context-awareness-manager/internal/database"
	"context-awareness-manager/internal/dockerclient"
	"context-awareness-manager/internal/monitor"
	"context-awareness-manager/internal/publisher"
	"context-awareness-manager/pkg/logger"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

// InitRouter initializes and returns the router with configured routes
func InitRouter(zenohURL string) (*mux.Router, monitor.ContextMonitor, *database.SQLConnector) {
	// Initialize logger
	logger.InitLogger()
	logrus.Info("Starting Context Awareness Manager...")

	// Initialize the database connection
	dbConnection, err := database.NewSQLConnector("./context_awareness_manager.db")
	if err != nil {
		logrus.Fatalf("Error initializing database: %v", err)
	}

	// Initialize Docker client
	dockerClient, err := dockerclient.NewDockerConnector()
	if err != nil {
		logrus.Fatalf("Error initializing Docker client: %v", err)
	}

	// Initialize the publisher
	publisher := publisher.NewPublisher(zenohURL)

	// Initialize the context service
	contextMonitor := monitor.NewContextMonitor(dbConnection, dockerClient, publisher)
	// Initialize the context controller
	contextHandler := controllers.NewContextHandler(contextMonitor)

	router := mux.NewRouter()
	initializeRoutes(router, contextHandler)
	return router, contextMonitor, dbConnection
}

// initializeRoutes initializes the routes for the server
func initializeRoutes(router *mux.Router, contextHandler *controllers.ContextHandler) {
	// Health check route
	router.HandleFunc("/healthcheck", controllers.HealthHandler).Methods("GET")

	// --- CRUD Context routes ---
	router.HandleFunc("/context", contextHandler.HandleContextRequest).Methods("POST")

	// Read all
	router.HandleFunc("/context", contextHandler.GetContexts).Methods("GET")

	// Read by ID
	router.HandleFunc("/context/{id}", contextHandler.GetContextByID).Methods("GET")

	// Delete by ID
	router.HandleFunc("/context/{id}", contextHandler.DeleteContext).Methods("DELETE")

	// Swagger documentation route
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	logrus.Info("Routes initialized successfully!")
}
