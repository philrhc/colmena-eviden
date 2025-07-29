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

package assessment

import (
	amodel "colmena/sla-management-svc/app/assessment/model"
	"colmena/sla-management-svc/app/common/logs"
	"colmena/sla-management-svc/app/model"
)

func GenerateViolationOutput(qos model.SLA, result amodel.Result) model.ColmenaOutputSLA {

	vs := result.GetViolations()
	if len(vs) == 0 {
		return model.ColmenaOutputSLA{}
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
	output, err := model.SLAModelToColmenaOutputSLA(qos)
	if err != nil {
		logs.GetLogger().Error("Error generating violation output: ", err)
		return model.ColmenaOutputSLA{}
	}

	return output
}
