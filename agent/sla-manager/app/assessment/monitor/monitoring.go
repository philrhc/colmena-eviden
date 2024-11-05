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
Package monitor contains the interface for monitoring adapters; subpackages contain several examples of adapters.
*/
package monitor

import (
	assessment_model "colmena/sla-management-svc/app/assessment/model"
	"colmena/sla-management-svc/app/model"

	"time"
)

// MonitoringAdapter is an interface which should be implemented per monitoring solution
type MonitoringAdapter interface {
	// Intialize the monitoring retrieval for one evaluation of the agreement
	//
	// A new MonitoringAdapter, copy of current adapter, must be returned
	Initialize(a *model.SLA) MonitoringAdapter

	// GetValues retrieve the metrics corresponding to the variables found in a guarantee
	GetValues(gt model.Guarantee, vars []string, to time.Time) assessment_model.GuaranteeData

	// Queries
	Query(metric string, path string) (interface{}, error)
}

// RetrievalItem contains the retrieval information for a variable
//
// Used in EarlyRetriever interface
type RetrievalItem struct {
	Guarantee model.Guarantee
	Var       model.Variable
	From      time.Time
	To        time.Time
}

// EarlyRetriever is implemented by adapters that want to (and can) retrieve
// all monitoring information in one query for efficiency reasons
type EarlyRetriever interface {
	RetrieveAllValues(items []RetrievalItem) []assessment_model.GuaranteeData
}
