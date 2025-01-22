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
// Package models contains the data structures used for service description, alerts, and KPIs.
package models

// ID represents the ID structure in the JSON
type ID struct {
	Value string `json:"value"`
}

// DockerContextDefinition represents a Docker context definition in the JSON
type DockerContextDefinition struct {
	ID      string `json:"id"`
	ImageID string `json:"imageId"`
}

// DockerRoleDefinition represents a Docker role definition in the JSON
type DockerRoleDefinition struct {
	ID                   string   `json:"id"`
	ImageID              string   `json:"imageId"`
	HardwareRequirements []string `json:"hardwareRequirements"`
	KPIs                 []string `json:"kpis"`
}

// ServiceDescription represents the entire JSON structure
type ServiceDescription struct {
	ID                       ID                        `json:"id"`
	DockerContextDefinitions []DockerContextDefinition `json:"dockerContextDefinitions"`
	KPIs                     []string                  `json:"kpis"`
	DockerRoleDefinitions    []DockerRoleDefinition    `json:"dockerRoleDefinitions"`
}

// Subscriber represents a subscriber with ID and endpoint
type Subscriber struct {
	ID       string `json:"id"`
	Endpoint string `json:"endpoint"`
}

type DeployRequest struct {
	Image string   `json:"image"`
	Cmd   []string `json:"cmd,omitempty"`
}

type DeployResponse struct {
	Classification string `json:"classification"`
	Error          string `json:"error,omitempty"`
}

// Result contiene el ID del contexto y la clasificación
type Result struct {
	ID             string `json:"id"`
	Classification string `json:"classification"`
}
