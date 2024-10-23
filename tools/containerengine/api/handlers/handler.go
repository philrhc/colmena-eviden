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
	"encoding/json"
	"net/http"

	"containerengine/internal/dockerclient"
	"containerengine/internal/models"

	"github.com/docker/docker/client"
)

// DeployHandler handles the deployment of a Docker container
// @Summary Deploy a Docker container
// @Description Deploys a container using the specified image and command
// @Tags deploy
// @Accept json
// @Produce json
// @Param deployRequest body models.DeployRequest true "Request payload for deploying a container"
// @Success 200 {object} models.DeployResponse
// @Failure 400 {object} models.DeployResponse
// @Router /deploy [post]
func DeployHandler(cli *client.Client, w http.ResponseWriter, r *http.Request) {
	var req models.DeployRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	logs, err := dockerclient.RunContainer(cli, req.Image, req.Cmd)
	response := models.DeployResponse{Classification: logs}
	if err != nil {
		response.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
