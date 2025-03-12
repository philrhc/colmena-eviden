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
package testadapter

import (
	"colmena/sla-management-svc/app/assessment/monitor"
	"colmena/sla-management-svc/app/assessment/monitor/genericadapter"
	"colmena/sla-management-svc/app/model"

	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"colmena/sla-management-svc/app/common/logs"

	"github.com/spf13/viper"
)

// path used in logs
const pathLOG string = "SLA-Framework > Assessment > Monitor > Test Adapter "

const (
	// Name is the unique identifier of this adapter/retriever
	Name = "testadapter"
)

// Retriever implements genericadapter.Retrieve
type Retriever struct{}

/*
New constructs a Prometheus adapter from a Viper configuration
*/
func New(config *viper.Viper) Retriever {
	logConfig(config)

	return Retriever{}
}

// logConfig
func logConfig(config *viper.Viper) {
	logs.GetLogger().Infof(pathLOG + " Configuration loaded.")
}

/*
Retrieve implements genericadapter.Retrieve
*/
func (r Retriever) Retrieve() genericadapter.Retrieve {
	return func(agreement model.SLA, items []monitor.RetrievalItem) map[model.Variable][]model.MetricValue {
		logs.GetLogger().Info(pathLOG + "[Retrieve] Retrieving metrics from Monitoring-Test Adapter ...")

		result := make(map[model.Variable][]model.MetricValue)
		for _, item := range items {
			logs.GetLogger().Info(pathLOG + "[Retrieve] Checking [item.Var.Name=" + item.Var.Name + "] ...")

			// call to test engine
			query := generateRandomResult(item.Var.Name)
			aux := translateVector(query, item.Var.Name)
			result[item.Var] = aux
		}
		logs.GetLogger().Infof(pathLOG+" Returning result: %v", result)

		return result
	}
}

// generateRandomResult
func generateRandomResult(metric string) query {
	data := query{}

	min := 10
	max := 1000
	v := rand.Intn(max-min) + min

	s := string(`{
		"status": "success",
		"data": {
			"resultType": "vector",
			"result": [
				{
					"metric": {
						"__name__": "` + metric + `",
						"instance": "localhost:9090",
						"job": "prometheus"
					},
					"value": [
						1571987564.298,
						"` + strconv.Itoa(v) + `"
					]
				}
			]
		}
	}`)

	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		logs.GetLogger().Error(pathLOG+"Error parsing random data to query struct: ", err.Error())
	}

	return data
}

// translateVector
func translateVector(query query, key string) []model.MetricValue {
	res := make([]model.MetricValue, 0, len(query.Data.Results))
	for _, item := range query.Data.Results {
		metric := translateMetric(key, item)
		res = append(res, metric)
	}
	return res
}

// translateMetric
// TODO this function should be made project-dependent
func translateMetric(key string, item result) model.MetricValue {
	return model.MetricValue{
		Key:      fmt.Sprintf("%s{%s}", key, item.Metric.Instance),
		Value:    item.Item.Value,
		DateTime: time.Time(item.Item.Timestamp),
	}
}
