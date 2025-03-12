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

// ServiceDescription represents the entire JSON structure
type ServiceDescription struct {
	ID                       interface{}               `json:"id"`
	DockerContextDefinitions []DockerContextDefinition `json:"dockerContextDefinitions"`
	KPIs                     interface{}               `json:"kpis"`
	DockerRoleDefinitions    interface{}               `json:"dockerRoleDefinitions"`
}

// DockerContextDefinition represents a Docker context definition in the JSON
type DockerContextDefinition struct {
	ID      string `json:"id"`
	ImageID string `json:"imageId"`
}

// DockerContextClassification contains the context ID and its classification.
type DockerContextClassification struct {
	ID             string `json:"id"`
	Classification string `json:"classification"`
}
