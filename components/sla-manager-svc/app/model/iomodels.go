/*
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

package model

const DEFAULT_ASSESSMENT_X = 2
const DEFAULT_ASSESSMENT_Y = 2
const DEFAULT_ASSESSMENT_Z = 5

/*
Service definition (input model example):

	{
		"id": {
			"value": "ExampleApplication"
		},
		"dockerContextDefinitions": [
			{
				"id": "company_premises",
				"imageId": "xaviercasasbsc/company_premises"
			}
		],
		"kpis": [],
		"dockerRoleDefinitions": [
			{
				"id": "test01",
				"imageId": "",
				"hardwareRequirements": [],
				"kpis": [{
					"query": "avg_over_time(processing_time[5s]) < 1",
					"scope": "company_premises/building=."
				}]
			},
			{
				"id": "test03",
				"imageId": "",
				"hardwareRequirements": [],
				"kpis": [{
					"query": "[sum%20by%20(metric_name, label1)%20(colmena_total_people{metric_name='tests', label1='planta01'})] < 5",
					"scope": ""
					}]
			},
			{
				"id": "test04",
				"imageId": "",
				"hardwareRequirements": [],
				"kpis": [{
					"query": "[sum by (metric_name, label1) (colmena_total_people{metric_name='tests', label1='planta01'})] < 5",
					"scope": ""
					}]
			}
		]
	}
*/
type InputSLA struct {
	ServiceId                ServiceId         `json:"id"`
	Roles                    []InputSLARole    `json:"dockerRoleDefinitions,omitempty"`
	DockerContextDefinitions []interface{}     `json:"dockerContextDefinitions,omitempty"`
	Kpis                     []InputSLARoleKPI `json:"kpis,omitempty"`
}

type ServiceId struct {
	Value string `json:"value"`
}

type InputSLARole struct {
	Id                   string            `json:"id,omitempty"`
	Kpis                 []InputSLARoleKPI `json:"kpis,omitempty"`
	ImageId              string            `json:"imageId,omitempty"`
	HardwareRequirements []interface{}     `json:"hardwareRequirements,omitempty"`
}

type InputSLARoleKPI struct {
	Query string `json:"query,omitempty"`
	Scope string `json:"scope,omitempty"`
}

/*
Output model (SLA VIOLATION) example:

	{
		"serviceId": "test_service_id",
		"KPIs": [
			{
				"roleId": "",
				"query": "avg(processing_time)<5",
				"value": 10,
				"threshold": 5,
				"level": "",
				"violations": [],
				"total_violations": 1
			}
		]
	}
*/
type OutputSLA struct {
	ServiceId string         `json:"serviceId"`
	SLAId     string         `json:"slaId"`
	Kpis      []OutputSLAKpi `json:"KPIs"`
}

type OutputSLAKpi struct {
	RoleId          string      `json:"roleId"`
	Query           string      `json:"query"`
	Level           string      `json:"level"`
	Value           interface{} `json:"value"`
	Threshold       float64     `json:"threshold"`
	Violations      []Violation `json:"violations,omitempty"`
	TotalViolations int         `json:"total_violations"`
}
