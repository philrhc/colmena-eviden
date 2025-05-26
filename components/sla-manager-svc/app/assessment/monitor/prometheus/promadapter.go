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
	"colmena/sla-management-svc/app/assessment/monitor"
	"colmena/sla-management-svc/app/assessment/monitor/genericadapter"
	"colmena/sla-management-svc/app/common/logs"
	"colmena/sla-management-svc/app/model"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

// path used in logs
const pathLOG string = "SLA > Assessment > Monitor > PROMETHEUS "

const (
	// Name is the unique identifier of this adapter/retriever
	Name = "prometheus"

	// PrometheusURLPropertyName is the config property name of the Prometheus URL
	PrometheusURLPropertyName = "PROMETHEUS_ADDRESS"

	// defaultURL is the value of the Prometheus URL is PrometheusURLPropertyName is not set
	defaultURL = "http://localhost:9090"
)

// Retriever implements genericadapter.Retrieve
type Retriever struct {
	URL string
}

/*
New constructs a Prometheus adapter from a Viper configuration
*/
func New(config *viper.Viper) Retriever {
	if os.Getenv(PrometheusURLPropertyName) != "" {
		config.Set(PrometheusURLPropertyName, os.Getenv(PrometheusURLPropertyName))
	} else {
		config.SetDefault(PrometheusURLPropertyName, defaultURL)
	}

	r := Retriever{
		config.GetString(PrometheusURLPropertyName),
	}

	logConfig(config)

	return r
}

// logConfig
func logConfig(config *viper.Viper) {
	logs.GetLogger().Info(pathLOG + "Prometheus configuration:\n" +
		"\t-----------------------------------------------------------------\n" +
		"\tURL (Prometheus location): " + config.GetString(PrometheusURLPropertyName) + "\n" +
		"\t-----------------------------------------------------------------")
}

/*
Retrieve implements genericadapter.Retrieve
*/
func (r Retriever) Retrieve() genericadapter.Retrieve {
	return func(agreement model.SLA, items []monitor.RetrievalItem) map[model.Variable][]model.MetricValue {
		rootURL := r.URL
		logs.GetLogger().Info(pathLOG + "[Retrieve] Retrieving metrics from Monitoring-PROMETHEUS adapter [" + rootURL + "] ...")

		result := make(map[model.Variable][]model.MetricValue)
		for _, item := range items {
			logs.GetLogger().Info(pathLOG + "[Retrieve] Checking [item.Var.Metric=" + item.Var.Metric + "], [item.Var.Name=" + item.Var.Name + "] ...")

			// call to monitoring engine:
			q := PromQLQuery{
				Metric: item.Var.Metric,
				Params: map[string]string{}}

			res := make([]model.MetricValue, 0, 1)
			for _, resQuery := range PromQuery(q.String()) {
				fv, err := strconv.ParseFloat(resQuery.Value.String(), 8)

				if err != nil {
					logs.GetLogger().Error(pathLOG + "ParseFloat Error: " + err.Error())
				} else {
					metric := model.MetricValue{
						Key:      item.Var.Name,
						Value:    fv,
						DateTime: resQuery.Timestamp.Time(),
					}
					res = append(res, metric)
				}

			}

			result[item.Var] = res
		}
		logs.GetLogger().Infof(pathLOG+" Returning result: ", result)

		return result
	}
}
