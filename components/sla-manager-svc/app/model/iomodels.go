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

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/lithammer/shortuuid/v4"
)

/**
	Input model example:

	{
    	"serviceId": "ExamplePlantcare",
    	"kpis": [],
    	"roles": [
    	{
        	"id": "Plantsensor",
        	"kpis": []
    	},
    	{
        	"id": "Plantwatering",
        	"kpis": ["exampleplantcare/moisture[5s] > 20"]
    	}]
	}

**/

type InputSLA struct {
	ServiceId string         `json:"serviceId"`
	Kpis      []string       `json:"kpis,omitempty"`
	Roles     []InputSLARole `json:"roles,omitempty"`
}

type InputSLARole struct {
	Id   string   `json:"id"`
	Kpis []string `json:"kpis,omitempty"`
}

/**
	Output model (SLA VIOLATION) example:

	{
  		"serviceId": null,
  		"roleId": "Plantwatering"
	}

**/

type OutputSLA struct {
	ServiceId string `json:"serviceId"`
	RoleId    string `json:"roleId"`
}

///////////////////////////////////////////////////////////////////////////////

/**
 * Transforms the input to an SLA Model
 */
func InputSLAModelToSLAModel(c *gin.Context) ([]SLA, error) {
	var input InputSLA
	var slas []SLA

	err := c.ShouldBindJSON(&input)
	if err != nil {
		return slas, err
	}

	// InputSLA ==> SLA(s)
	if len(input.Roles) > 0 {
		for _, r := range input.Roles {
			if len(r.Kpis) > 0 {
				uid := uuid.New()
				sla := SLA{}

				sla.Name = input.ServiceId
				sla.Id = input.ServiceId + "-" + uid
				sla.State = "started"

				sla.Details.Guarantees = make([]Guarantee, 1) // TODO for each KPI => 1 Guarantee
				sla.Details.Guarantees[0].Name = r.Id
				sla.Details.Guarantees[0].Constraint = r.Kpis[0]

				slas = append(slas, sla)
			}
		}
	}

	return slas, nil
}
