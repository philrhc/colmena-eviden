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
	cfgconst "colmena/sla-management-svc/app/common/cfg"
	"colmena/sla-management-svc/app/common/expressions"
	"colmena/sla-management-svc/app/common/logs"
	"colmena/sla-management-svc/app/model"
	"encoding/json"
	"io"
	"reflect"
	"strings"

	"net/http"

	"github.com/spf13/viper"
)

// Context Definitions struct response
/*
ResponseData example:
	{
		"key":"colmena/contexts/ColmenaAgent1/company_premises",
		"value":{"building":"Red","floor":"22","room":"Rest Room"},
		"encoding":"application/json",
		"timestamp":""
	}
*/
type ResponseData struct {
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
	Encoding  string      `json:"encoding,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
}

/*
CheckPausedQoSDefinitions will check the paused QoS from the provided repository and set them to active when the context
is completed. PAUSED SLA example:

	{
		"id": "ExampleApplication-b7MJb7imo5qqcA4Ah25qRC",
		"name": "ExampleApplication",
		"state": "paused",
		"assessment": {
			"first_execution": "0001-01-01T00:00:00Z",
			"last_execution": "0001-01-01T00:00:00Z"
		},
		"creation": "2025-04-02T09:50:13.2744673+01:00",
		"expiration": "2026-04-02T09:50:13.2744673+01:00",
		"details": {
			"guarantees": [{
					"name": "Processing",
					"constraint": "[go_memstats_frees_total] < 50000",
					"query": "[go_memstats_frees_total#LABELS#] < 50000",
					"scope": "company_premises/building=.",
					"scopeTemplate": "company_premises/building=."
				}
			]
		}
	}
*/
func CheckPausedQoSDefinitions(cfg Config, vconfig *viper.Viper) {
	repo := cfg.Repo

	// Retrieve all PAUSED SLAs
	qosdefs, err := repo.GetSLAsByState(model.PAUSED)
	if err != nil {
		logs.GetLogger().Error(pathLOG+"[CheckPausedQoSDefinitions] Error getting PAUSED SLAs: %s", err.Error())
	} else {
		logs.GetLogger().Infof(pathLOG+"[CheckPausedQoSDefinitions] [%d PAUSED SLAs to check]", len(qosdefs))

		if len(qosdefs) > 0 {
			// get results from Zenoh
			var items []ResponseData = getContextResults(vconfig)
			if len(items) > 0 {
				// check PAUSED SLAs
				for _, qosd := range qosdefs {
					logs.GetLogger().Info(pathLOG + "[CheckPausedQoSDefinitions] Checking PAUSED SLA " + qosd.Id + " ...")

					// if context updated
					if checkSLA(&qosd, items, vconfig) {
						// update SLA (context and status)
						repo.UpdateSLA(&qosd)
						logs.GetLogger().Info(pathLOG + "[CheckPausedQoSDefinitions] SLA " + qosd.Id + " set to STARTED ...")
					}
				}
			}
		}
	}
}

/*
getContextResults does a query to Zenoh (e.g. GET http://192.168.137.47:8000/colmena/contexts/**) to get all values
from context. Example:

curl http://192.168.137.47:8000/colmena/contexts/**

	[
		{
			"key":"colmena/contexts/ColmenaAgent1/company_premises",
			"value":{"building":"Red","floor":"22","room":"Rest Room"},
			"encoding":"application/json",
			"timestamp":""
		}
		...
	]

Returns a list of ResponseData objects.
*/
func getContextResults(vconfig *viper.Viper) []ResponseData {
	res := []ResponseData{}

	endpoint := vconfig.GetString(cfgconst.ContextZenohEndpointPropertyName) +
		vconfig.GetString(cfgconst.ContextZenohContextsPropertyName) + "/**"
	logs.GetLogger().Info(pathLOG + "[getContextResults] Getting contexts from Zenoh endpoint [" + endpoint + "] ...")
	resp, err := http.Get(endpoint)

	if err != nil {
		logs.GetLogger().Error(pathLOG+"Error:", err)
	} else {
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logs.GetLogger().Error(pathLOG+"Error:", err)
		} else {
			// Print the response body
			logs.GetLogger().Debug(pathLOG + " RESP: " + string(body))

			// Parse the JSON response
			var items []ResponseData
			err = json.Unmarshal(body, &items)
			if err != nil {
				logs.GetLogger().Error("Failed to parse JSON: %v", err)
			} else {
				// Print the parsed items
				for _, item := range items {

					resultMap := item.Value.(map[string]interface{})
					for key, value := range resultMap {
						logs.GetLogger().Debug(pathLOG+" Key: ", key, " Value: ", value, " Type: ", reflect.TypeOf(value))
					}

					//logs.GetLogger().Debug(pathLOG+"Key: %s, Value: %s\n", item.Key, resultMap["building"].(string))
					res = append(res, item)
				}
			}
		}
	}

	return res
}

/*
SLA

	{
		...
		"details": {
			"guarantees": [{
					"name": "Processing",
					"constraint": "[go_memstats_frees_total] < 50000",
					"query": "[go_memstats_frees_total#LABELS#] < 50000",
					"scope": "company_premises/building=.",
					"scopeTemplate": "company_premises/building=."
				}
			]
		}
	}

[]ResponseData

	[
		{
			"key":"colmena/contexts/ColmenaAgent1/company_premises",
			"value":{"building":"Red","floor":"22","room":"Rest Room"},
			"encoding":"application/json",
			"timestamp":""
		}
		...
	]
*/
func checkSLA(sla *model.SLA, items []ResponseData, vconfig *viper.Viper) bool {
	logs.GetLogger().Debug(pathLOG + "[checkSLA] Checking SLA ...")

	item := getData(*sla, items, vconfig)
	/*
		item example:
			{
				"key":"colmena/contexts/ColmenaAgent1/company_premises",
				"value":{"building":"Red","floor":"22","room":"Rest Room"},
				"encoding":"application/json",
				"timestamp":""
			}
	*/
	if (item == ResponseData{}) {
		logs.GetLogger().Debug(pathLOG + "[checkSLA] No context for SLA found.")
		return false
	}

	/*

		{
			"key":"colmena/contexts/ColmenaAgent1/company_premises",
			"value":{"building":"Red","floor":"22","room":"Rest Room"},
			"encoding":"application/json",
			"timestamp":""
		}

		==>

		LABEL_MARK = {company_premises_building="\Red"\}

	*/

	fullContextLabel := getFullContextLabelFromScope(*sla) // e.g. "company_premises_building"
	destLabel := getSearchLabelFromScope(*sla)             // e.g. "building"
	destLabelValue := ""

	resultMap := item.Value.(map[string]interface{})
	for key, value := range resultMap {
		if key == destLabel {
			logs.GetLogger().Debug(pathLOG+"Label from scope: ", destLabel, " Key: ", key, " Value: ", value, " Type: ", reflect.TypeOf(value))

			str, ok := value.(string)
			if ok {
				destLabelValue = str
			} else {
				logs.GetLogger().Warn(pathLOG+" Key: ", key, " Value: ", value, " Type: ", reflect.TypeOf(value))
			}
		}
	}

	if len(fullContextLabel) > 0 && len(destLabel) > 0 && len(destLabelValue) > 0 {
		sla.State = model.STARTED
		// replace labels
		q := strings.Replace(sla.Details.Guarantees[0].Query,
			expressions.LABEL_MARK,
			"{"+fullContextLabel+"=\""+destLabelValue+"\"}",
			1)
		sla.Details.Guarantees[0].Constraint = q

		return true
	}

	logs.GetLogger().Debug(pathLOG + "[checkSLA] Label used in PromQL queries not set.")
	return false
}

/*
SLA example:

	{
		...
		"details": {
			"guarantees": [{
					...
					"scope": "company_premises/building=."
				}
			]
		}
	}

[]ResponseData

	[
		{
			"key":"colmena/contexts/ColmenaAgent1/company_premises",
			"value":{"building":"Red","floor":"22","room":"Rest Room"},
			"encoding":"application/json",
			"timestamp":""
		}
		...
	]
*/
func getData(sla model.SLA, items []ResponseData, vconfig *viper.Viper) ResponseData {
	context_scope := getSLAContextFromScope(sla) // e.g. "company_premises"
	if len(context_scope) > 0 {
		context := ""
		if len(vconfig.GetString(cfgconst.AgentIdPropertyName)) > 0 {
			context = vconfig.GetString(cfgconst.ContextZenohContextsPropertyName) + "/" +
				vconfig.GetString(cfgconst.AgentIdPropertyName) + "/" + context_scope
		} else {
			context = vconfig.GetString(cfgconst.ContextZenohContextsPropertyName) + "/" +
				vconfig.GetString(cfgconst.ComposeProjectPropertyName) + "/" + context_scope
		}

		// => e.g. "colmena/contexts/ColmenaAgent1/company_premises"
		logs.GetLogger().Debug(pathLOG + "[checkSLA] context: " + context)

		for _, item := range items {
			if item.Key == context {
				return item
			}
		}
	}

	logs.GetLogger().Warn(pathLOG + "[getData] Context data not found")
	return ResponseData{}
}

/*
SLA example:

	{
		...
		"details": {
			"guarantees": [{
					...
					"scope": "company_premises/building=."
				}
			]
		}
	}

Returns "company_premises"
*/
func getSLAContextFromScope(sla model.SLA) string {
	logs.GetLogger().Debug(pathLOG + "[getSLAContextFromScope] Getting context from SLA scope ...")
	if len(sla.Details.Guarantees) > 0 {
		scope := sla.Details.Guarantees[0].Scope
		if len(scope) > 0 {
			arr := strings.Split(scope, "/")
			if len(arr) > 0 {
				logs.GetLogger().Debug(pathLOG + "[getSLAContextFromScope] return " + arr[0])
				return arr[0]
			}
		}
	}

	logs.GetLogger().Warn(pathLOG + "[getSLAContextFromScope] Error with SLA context. Return ''")
	return ""
}

/*
SLA example:

	{
		...
		"details": {
			"guarantees": [{
					...
					"scope": "company_premises/building=."
				}
			]
		}
	}

Returns "company_premises_building"
*/
func getFullContextLabelFromScope(sla model.SLA) string {
	logs.GetLogger().Debug(pathLOG + "[getFullContextLabelFromScope] Getting full context label from SLA scope ...")
	if len(sla.Details.Guarantees) > 0 {
		scope := sla.Details.Guarantees[0].Scope
		if len(scope) > 0 {
			arr := strings.Split(scope, "/")
			if len(arr) > 1 && len(arr[1]) > 2 {
				logs.GetLogger().Debug(pathLOG + "[getFullContextLabelFromScope] return " + arr[0] + "_" + arr[1][:len(arr[1])-2])
				return arr[0] + "_" + arr[1][:len(arr[1])-2]
			}
		}
	}

	logs.GetLogger().Warn(pathLOG + "[getFullContextLabelFromScope] Full context label error. Return ''")
	return ""
}

/*
SLA example:

	{
		...
		"details": {
			"guarantees": [{
					...
					"scope": "company_premises/building=."
				}
			]
		}
	}

Returns "building"
*/
func getSearchLabelFromScope(sla model.SLA) string {
	logs.GetLogger().Debug(pathLOG + "[getSearchLabelFromScope] Getting search label from SLA scope ...")
	if len(sla.Details.Guarantees) > 0 {
		scope := sla.Details.Guarantees[0].Scope
		if len(scope) > 0 {
			arr := strings.Split(scope, "/")
			if len(arr) > 1 && len(arr[1]) > 2 {
				logs.GetLogger().Debug(pathLOG + "[getSearchLabelFromScope] return '" + arr[1][:len(arr[1])-2] + "'")
				return arr[1][:len(arr[1])-2]
			}
		}
	}

	logs.GetLogger().Warn(pathLOG + "[getSearchLabelFromScope] Search label error. Return ''")
	return ""
}
