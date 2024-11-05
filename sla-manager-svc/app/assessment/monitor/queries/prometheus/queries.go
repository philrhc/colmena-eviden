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

	// defaultURL is the value of the Prometheus URL is PrometheusURLPropertyName is not set
	defaultURL = "http://localhost:9090"
)

type PromQLQuery struct {
	Metric string
	Params map[string]string
}

// Query implements monitor.MonitoringAdapter.Query
func Query(metric string, path string) (interface{}, error) {
	logs.GetLogger().Debug("metric: " + metric + ", path: " + path)

	// call to monitoring engine:
	q := PromQLQuery{
		Metric: "metric",
		Params: map[string]string{}}
	for _, resQuery := range PromQuery(q.String()) {
		logs.GetLogger().Debug("resQuery: " + resQuery.Value.String())
	}

	logs.GetLogger().Warn("Function not implemented")
	return nil, nil
}

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

func PromQuery(query string) model.Vector {

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
