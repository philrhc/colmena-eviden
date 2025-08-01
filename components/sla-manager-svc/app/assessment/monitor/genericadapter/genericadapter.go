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

/*
Package genericadapter provides a configurable MonitoringAdapter
that works with advanced agreement schema.

Usage:
	ma := genericadapter.New(retriever, processor)
	ma = ma.Initialize(&agreement)
	for _, gt := range gts {
		for values := range ma.GetValues(gt, ...) {
			...
		}
	}
*/

package genericadapter

import (
	"colmena/sla-management-svc/app/assessment"
	amodel "colmena/sla-management-svc/app/assessment/model"
	monitor "colmena/sla-management-svc/app/assessment/monitor"
	query_prometheus "colmena/sla-management-svc/app/assessment/monitor/queries/prometheus"
	"colmena/sla-management-svc/app/common/logs"
	"colmena/sla-management-svc/app/model"
	"strings"

	"math/rand"
	"time"
)

/*
Adapter is the type of a customizable adapter.

The Retrieve field is a function to query data to monitoring;
the Process field is a function to perform additional processing
on data.

Two Process functions are provided in the package:
Identity (returns the input) and Aggregation (aggregates values according
to the aggregation type)
*/
type Adapter struct {
	Type      string
	Retrieve  Retrieve
	Process   Process
	agreement *model.SLA
}

// Retrieve is the type of the function that makes the actual request to monitoring.
//
// It receives the list of variables to be able to retrieve all of them at once if possible.
type Retrieve func(agreement model.SLA, items []monitor.RetrievalItem) map[model.Variable][]model.MetricValue

// Process is the type of the function that performs additional custom processing on
// retrieved data.
type Process func(v model.Variable, values []model.MetricValue) []model.MetricValue

// New is a helper function to build an Adapter from a Retriever and the Process function.
func New(t string, retrieve Retrieve, process Process) monitor.MonitoringAdapter {
	return &Adapter{
		Type:     t,
		Retrieve: retrieve,
		Process:  process,
	}
}

/*
Initialize implements MonitoringAdapter.Initialize().

Usage:

	ga := GenericAdapter{
		Retrieve: randomRetrieve,
		Process: Aggregation,
	}
	ga := ga.Initialize(agreement)
	for _, gt := range gts {
		for values := range ga.GetValues(gt, ...) {
			...
		}
	}
*/
func (ga *Adapter) Initialize(a *model.SLA) monitor.MonitoringAdapter {
	result := *ga
	result.agreement = a
	return &result
}

// Query realizes a query to monitoring adapter (i.e. Prometheus) to get the metric values
func (ga *Adapter) Query(metric string, path string) (interface{}, error) {
	logs.GetLogger().Debug("Adapter: " + ga.Type)
	if strings.ToLower(ga.Type) == strings.ToLower("prometheus") {
		return query_prometheus.Query(metric, path)
	}

	logs.GetLogger().Warn("Query adapter not supported")
	return nil, nil
}

// GetValues implements Monitoring.GetValues().
func (ga *Adapter) GetValues(gt model.Guarantee, varnames []string, now time.Time) amodel.GuaranteeData {

	a := ga.agreement

	items := assessment.BuildRetrievalItems(a, gt, varnames, now)
	unprocessed := ga.Retrieve(*a, items)

	/* process each of the series*/
	valuesmap := map[model.Variable][]model.MetricValue{}
	for v := range unprocessed {
		valuesmap[v] = ga.Process(v, unprocessed[v])
	}
	result := Mount(valuesmap, lastvalues(a, gt), 0.1)
	return result
}

func lastvalues(a *model.SLA, gt model.Guarantee) model.LastValues {
	empty := model.LastValues{}
	if a.Assessment.Guarantees == nil {
		return empty
	}
	ag, ok := a.Assessment.Guarantees[gt.Name]
	if !ok {
		return empty
	}
	if ag.LastValues == nil {
		return empty
	}
	return ag.LastValues
}

// DummyRetriever is a simple struct that generates a RetrieveFunction that works similar
// to the DummyAdapter, returning random values for each variable.
//
// Usage:
//
//	adapter := Adapter { Retrieve: DummyRetriever{3}.RetrieveFunction() }
type DummyRetriever struct {
	// Size is the number of values that the retrieval returns per metric
	Size int
}

// Retrieve returns a Retrieve function.
func (r DummyRetriever) Retrieve() Retrieve {

	return func(agreement model.SLA,
		items []monitor.RetrievalItem) map[model.Variable][]model.MetricValue {

		result := map[model.Variable][]model.MetricValue{}
		for _, item := range items {
			v := item.Var
			actualFrom := item.From
			to := item.To
			result[v] = make([]model.MetricValue, 0, r.Size)
			step := time.Duration(int(to.Sub(actualFrom)) / (r.Size + 1))

			for i := 0; i < r.Size; i++ {
				m := model.MetricValue{
					Key:      v.Name,
					Value:    rand.Float64(),
					DateTime: actualFrom.Add(step * time.Duration(i+1)),
				}
				result[v] = append(result[v], m)
			}
		}
		return result
	}
}

// Identity returns the input
func Identity(v model.Variable, values []model.MetricValue) []model.MetricValue {
	return values
}

// Aggregate performs an aggregation function on the input.
//
// This expects that all the values are in the appropriate window. For that,
// the Retrieve function needs to return only the values in the window. If not,
// this function will return an invalid result.
func Aggregate(v model.Variable, values []model.MetricValue) []model.MetricValue {
	if len(values) == 0 || v.Aggregation == nil || v.Aggregation.Type == "" {
		return values
	}
	if v.Aggregation.Type == model.AVERAGE {
		avg := average(values)
		return []model.MetricValue{
			model.MetricValue{
				Key:      v.Name,
				Value:    avg,
				DateTime: values[len(values)-1].DateTime,
			},
		}
	}
	/* fallback */
	return values
}

func average(values []model.MetricValue) float64 {
	sum := 0.0
	for _, value := range values {
		sum += value.Value.(float64)
	}
	result := sum / float64(len(values))

	return result
}
