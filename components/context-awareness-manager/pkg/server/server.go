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
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

// StartServer initializes the server and handles all HTTP requests.
func StartServer(port string, router *mux.Router) {
	// Create the HTTP server instance
	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           cors.AllowAll().Handler(router),
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Capture shutdown signals and perform graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdownChan)

	// Start the server in a goroutine
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Server failed to start: %v", err)
		}
	}()
	logrus.Info("Server is ready to handle requests")

	// Wait for shutdown signal
	<-shutdownChan
	gracefulShutdown(httpServer)
}

// gracefulShutdown ensures the server stops cleanly
func gracefulShutdown(httpServer *http.Server) {
	logrus.Info("Shutdown signal received, stopping server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logrus.Errorf("Server Shutdown Failed: %+v", err)
	} else {
		logrus.Info("Server gracefully stopped")
	}
}
