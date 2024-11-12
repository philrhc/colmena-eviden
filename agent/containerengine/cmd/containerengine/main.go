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
	"fmt"
	"log"
	"net/http"
	"os"

	"containerengine/api/handlers"

	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

func main() {

	host := os.Getenv("DOCKER_HOST")
	if host == "" {
		host = "unix:///var/run/docker.sock"
	}
	// Create a Docker client
	cli, err := client.NewClientWithOpts(client.WithHost(host), client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("error creating Docker client: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
	r.HandleFunc("/deploy", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeployHandler(cli, w, r)
	}).Methods("POST")

	fmt.Println("Starting server on :9000")
	log.Fatal(http.ListenAndServe(":9000", r))
}
