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
package prometheus

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"colmena/sla-management-svc/app/common/logs"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// path used in logs
const pathLOG string = "SLA > Assessment > Monitor > Queries "

const (
	// Name is the unique identifier of this adapter/query
	Name = "prometheus"

	// PrometheusURLPropertyName is the config property name of the Prometheus URL
	PrometheusURLPropertyName = "PROMETHEUS_ADDRESS"
)

type PromQLQuery struct {
	Metric string
	Params map[string]string
}

type PromQLQueryResponse struct {
	Metrics    []model.Sample
	Aggregated PromQLQueryAggregated
}

type PromQLQueryAggregated struct {
	Total string
}

// Query implements monitor.MonitoringAdapter.Query
func Query(metric string, path string) (interface{}, error) {
	logs.GetLogger().Debug("metric: " + metric + ", path: " + path)

	params := make(map[string]string)
	params["path"] = path

	// call to monitoring engine:
	// Prometheus queries examples:
	// Get metric values
	// 	- http://192.168.137.47:9090/api/v1/query?query=colmena_metric1
	// 	- metric = 'colmena_metric1'
	// Get metric values filtered by label
	// 	- http://192.168.137.47:9090/api/v1/query?query=colmena_metric1{path=%22/tests/planta01/habitacion01%22}
	// 	- metric = 'colmena_metric1'
	// 	- path = '/tests/planta01/habitacion01'
	// Get metric value filtered by label using regex
	// 	- http://192.168.137.47:9090/api/v1/query?query=colmena_metric1{path=~%22/tests/planta01.*%22}
	// 	- metric = 'colmena_metric1'
	// 	- path = '~/tests/planta01.*'
	q := PromQLQuery{
		Metric: metric,
		Params: params}

	// by default, if path includes a regex expression '~',
	// 		- http://localhost:8080/api/v1/query?metric=colmena_metric1&path=~/tests/planta01/.*
	//		- path = "~/tests/planta01/.*"
	// the character "~" is included as part of the path value and not before the string
	// 		- path = ~"/tests/planta01/.*"
	//
	// Here the '~' is moved out of the path value
	query_string := q.String()
	logs.GetLogger().Debug("query_string [" + query_string + "]")

	if strings.Contains(path, "~") && strings.HasPrefix(path, "~") {
		logs.GetLogger().Debug("path value [" + path + "] is a (prometheus) regular expression.")
		query_string = strings.Replace(query_string, "~", "", 1)
		logs.GetLogger().Debug("query_string [" + query_string + "]")
		query_string = strings.Replace(query_string, "path=", "path=~", 1)
		logs.GetLogger().Debug("query_string [" + query_string + "]")
	}

	result := PromQLQueryResponse{}

	// Query Prometheus
	l := PromQuery(query_string)
	for _, resQuery := range l {
		result.Metrics = append(result.Metrics, *resQuery)
		logs.GetLogger().Debug("resQuery: " + resQuery.Value.String())
		logs.GetLogger().Debug("resQuery: " + resQuery.Metric.String())

		result.Aggregated.Total = CalcAggregated(resQuery.Value.String(), result.Aggregated.Total)
	}

	//return l, nil
	return result, nil
}

// CalcAggregated calculates the aggregated
func CalcAggregated(v string, i string) string {

	ir, _ := strconv.ParseInt(i, 10, 64)

	if r, err := strconv.ParseInt(v, 10, 64); err == nil {
		ir = ir + r
	}

	return strconv.Itoa(int(ir))
}

// PromQLQuery struct to string
func (q PromQLQuery) String() string {

	params := make([]string, 0, len(q.Params))
	for key, value := range q.Params {
		params = append(params, key+`="`+value+`"`)
	}
	if len(params) > 0 {
		return q.Metric + "{" + strings.Join(params, ", ") + "}"
	} else {
		return q.Metric
	}

}

// Query Prometheus
func PromQuery(query string) model.Vector {

	// HTML URL Encoding Reference https://www.w3schools.com/tags/ref_urlencode.asp
	// Examples:
	// 	~ 	%7E
	// 	[ 	%5B
	// 	]	%5D
	query = strings.ReplaceAll(query, "%5B", "[")
	query = strings.ReplaceAll(query, "%5D", "]")
	// create prometheus API client
	client, err := api.NewClient(api.Config{
		Address: os.Getenv("PROMETHEUS_ADDRESS"),
	})
	if err != nil {
		logs.GetLogger().Error(pathLOG+"Error creating client: ", err)
	}

	logs.GetLogger().Debug(pathLOG+"PromQL Query: ", query)

	// create prometheus API object
	v1api := v1.NewAPI(client)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, warnings, err := v1api.Query(ctx, query, time.Now(), v1.WithTimeout(10*time.Second))
	if err != nil {
		logs.GetLogger().Error(pathLOG+"Error querying Prometheus: ", err)
	}
	if len(warnings) > 0 {
		logs.GetLogger().Warn(pathLOG+"Warnings: ", warnings)
	}

	// match the response to vector and print the response values
	switch r := result.(type) {
	case model.Vector:

		if r.Len() == 0 {
			logs.GetLogger().Debug(pathLOG+"PromQL Query Result length is zero: ", query)
		}

		return r

	default:
		logs.GetLogger().Error(pathLOG + "Response is not a modelVector. Error: " + errors.New("not implemented").Error())

		return nil
	}
}
