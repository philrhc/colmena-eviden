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
	docker "containerengine/src/dockerclient"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type DeployRequest struct {
	Image string   `json:"image"`
	Cmd   []string `json:"cmd,omitempty"`
}

type DeployResponse struct {
	Classification string `json:"classification"`
	Error          string `json:"error,omitempty"`
}

func deployHandler(w http.ResponseWriter, r *http.Request) {
	var req DeployRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	logs, err := docker.RunContainer(req.Image, req.Cmd)
	response := DeployResponse{Classification: logs}
	if err != nil {
		response.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/deploy", deployHandler).Methods("POST")

	fmt.Println("Starting server on :9000")
	log.Fatal(http.ListenAndServe(":9000", r))
}
