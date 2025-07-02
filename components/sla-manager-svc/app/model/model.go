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
Package model contains the entities used in the SLA assessment: qos, violations, penalties etc.
It also defines the interface IRepository, which defines the operations to be implemented by any repository.
*/
package model

import (
	"errors"
	"fmt"
	"time"
)

/**
	SLA example:

	{
		"id": "ExampleApplication-XWBnySXE26VFnNcv429jn5",
		"name": "ExampleApplication",
		"state": "started",
		"total_executions": 1,
		"total_violations": 1,
		"assessment": {
			"first_execution": "2025-04-07T11:17:13.6064087+01:00",
			"last_execution": "2025-04-07T11:17:13.6064087+01:00",
			"guarantees": {
				"Processing": {
					"first_execution": "2025-04-07T11:17:13.6064087+01:00",
					"last_execution": "2025-04-07T11:17:13.6064087+01:00",
					"last_values": {
						"go_memstats_frees_total": {
							"key": "go_memstats_frees_total",
							"action": "",
							"namespace": "",
							"value": 514512460,
							"datetime": "2025-04-07T11:17:13.614+01:00"
						}
					},
					"last_violation": {
						"id": "",
						"agreement_id": "ExampleApplication-XWBnySXE26VFnNcv429jn5",
						"guarantee": "Processing",
						"action": "",
						"datetime": "2025-04-07T11:17:13.614+01:00",
						"constraint": "[go_memstats_frees_total] \u003C 50000",
						"values": [{
								"key": "go_memstats_frees_total",
								"action": "",
								"namespace": "",
								"value": 514512460,
								"datetime": "2025-04-07T11:17:13.614+01:00"
							}
						],
						"importanceName": "Default",
						"importance": -1,
						"appID": "ExampleApplication-XWBnySXE26VFnNcv429jn5"
					}
				}
			}
		},
		"creation": "2025-04-07T11:16:33.2662685+01:00",
		"expiration": "2026-04-07T11:16:33.2662685+01:00",
		"details": {
			"guarantees": [{
					"name": "Processing",
					"constraint": "[go_memstats_frees_total] \u003C 50000",
					"query": "[go_memstats_frees_total] \u003C 50000",
					"scope": "",
					"scopeTemplate": ""
				}
			]
		}
	}

**/

// ErrNotFound is the sentinel error for an entity not found
var ErrNotFound = errors.New("Entity not found")

// ErrAlreadyExist is the sentinel error for creating an entity whose id already exists
var ErrAlreadyExist = errors.New("Entity already exists")

/*
 * ValidationErrors following behavioral errors
 * (https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
 */

// validationError is an interface that must be implemented by custom error implementations
type validationError interface {
	IsErrValidation() bool
}

// IsErrValidation return true is an error is a validation error
func IsErrValidation(err error) bool {
	v, ok := err.(validationError)
	return ok && v.IsErrValidation()
}

// func IsErrNotFound(err error) bool

// Identity identifies entities with an Id field
type Identity interface {
	GetId() string
}

// Validable identifies entities that can be validated
type Validable interface {
	Validate(val Validator, mode ValidationMode) []error
}

// State is the type of possible states of an agreement
type State string

// TextType is the type of possible types a Details type
type TextType string

// AggregationType is the type of supported variable aggregations
type AggregationType string

const (
	// STARTED is the state of an agreement that can be evaluated
	STARTED State = "started"

	// STOPPED is the state of an agreement temporaryly not evaluated
	STOPPED State = "stopped"

	// TERMINATED is the final state of an agreement
	TERMINATED State = "terminated"

	// PAUSED is the state of an agreement temporaryly paused
	PAUSED State = "paused"

	// INVALID is the state of an agreement that is not valid
	INVALID State = "invalid"
)

const (
	// NONE is used when the variable is not aggregated
	NONE AggregationType = "none"
	// AVERAGE is used to calculate average of a variable
	AVERAGE AggregationType = "average"
)

// States is the list of possible states of an agreement/template
var States = [...]State{STOPPED, STARTED, TERMINATED}

///////////////////////////////////////////////////////////////////////////////
// SLA Model
/*
 Example:

	{
		"id": "d01",
		"name": "qos-definition-name",
		"state": "started",
		"assessment": {}
		"creation": "2024-01-16T17:09:45Z",
		"expiration": "2026-01-16T17:09:45Z",
		"details": {
			"variables": []
			"guarantees": []
		}
	}

*/

// SLA is the entity that represents a SLA definition.
// The Text is ReadOnly in normal conditions, with the exception of a renegotiation.
// The Assessment cannot be modified externally.
type SLA struct {
	Id         string     `json:"id" bson:"_id"`
	Name       string     `json:"name"`
	State      State      `json:"state"`
	Assessment Assessment `json:"assessment,omitempty"`
	Creation   time.Time  `json:"creation,omitempty"`
	Expiration *time.Time `json:"expiration,omitempty"`
	Details    Details    `json:"details"`
}

// Details is the struct that represents the "contract" signed by the client
type Details struct {
	Variables  []Variable  `json:"variables,omitempty"`
	Guarantees []Guarantee `json:"guarantees"`
}

// Variable gives additional information about a metric used in a Guarantee constraint
type Variable struct {
	Name        string       `json:"name"`
	Metric      string       `json:"metric"`
	Aggregation *Aggregation `json:"aggregation,omitempty"`
}

// Assessment is the struct that provides assessment information
type Assessment struct {
	TotalExecutions int `json:"total_executions,omitempty"` // total executions
	TotalViolations int `json:"total_violations,omitempty"` // total violations
	/*
		1) Primera vez que una KPI se infringe -> Level = Broken
		2) Después de X veces [seguidas] que se ha infringido -> Level = Critical
		3) Primera vez que una KPI se cumple [después de estar Broken] -> Level = Met
		4) Después de Y veces [seguidas] que se ha cumplido -> Level = Desired
		5) Sí ha cambiado de KPI met a KPI broken Z veces -> Level = Unstable

		TODO: remove from json => `json:"-"`
	*/
	X              int                            `json:"x,omitempty"`                         // assessment violation; x (default 2)
	XCounter       int                            `json:"x_assessment_broken_count,omitempty"` // assessment violation counter
	Y              int                            `json:"y,omitempty"`                         // assessment met; y (default 2)
	YCounter       int                            `json:"y_assessment_met_count,omitempty"`    // assessment met; y counter
	Z              int                            `json:"z,omitempty"`                         // met to broken; z (default 5)
	ZCounter       int                            `json:"z_met_to_broken_count,omitempty"`     // met to broken; z counter
	Level          string                         `json:"level,omitempty"`                     // Broken, Critical, Met, Desired, Unstable, Unknown
	Threshold      float64                        `json:"threshold,omitempty"`
	Violated       bool                           `json:"violated,omitempty"`
	FirstExecution time.Time                      `json:"first_execution"`
	LastExecution  time.Time                      `json:"last_execution"`
	MonitoringURL  string                         `json:"monitoring_url,omitempty"`
	Guarantees     map[string]AssessmentGuarantee `json:"guarantees,omitempty"` // Guarantees may be nil. Use Assessment.SetGuarantee to create if needed.
}

// Broken, Critical, Met, Desired, Unstable, Unknown
const (
	ASSESSMENT_LEVEL_UNKNOWN   = "Unknown"
	ASSESSMENT_LEVEL_NORESULTS = "Unknown_NoResults"
	ASSESSMENT_LEVEL_UNSTABLE  = "Unstable"
	ASSESSMENT_LEVEL_DESIRED   = "Desired"
	ASSESSMENT_LEVEL_MET       = "Met"
	ASSESSMENT_LEVEL_CRITICAL  = "Critical"
	ASSESSMENT_LEVEL_BROKEN    = "Broken"
)

// AssessmentGuarantee contain the assessment information for a guarantee term
type AssessmentGuarantee struct {
	FirstExecution time.Time  `json:"first_execution"`
	LastExecution  time.Time  `json:"last_execution"`
	LastValues     LastValues `json:"last_values,omitempty"`
	LastViolation  *Violation `json:"last_violation,omitempty"`
}

// LastValues contain last values of variables in guarantee terms
type LastValues map[string]MetricValue

// Guarantee is the struct that represents an SLO
type Guarantee struct {
	Name          string `json:"name"`
	Constraint    string `json:"constraint"`
	Query         string `json:"query"`
	OQuery        string `json:"oquery"`
	Scope         string `json:"scope"`
	ScopeTemplate string `json:"scopeTemplate"`
}

// Aggregation gives aggregation information of a variable.
// If defined and value is not NONE, the metric must be aggregated
// in the specified window in seconds.
// I.e. (average, 3600) means that the average over a period of one hour is calculated.
type Aggregation struct {
	Type   AggregationType `json:"type"`
	Window int             `json:"window"`
}

// For tracking reasons, we add the type of violation
// to define "mild", "serious", "catastrophic", etc. violations
type GuaranteeType struct {
	Name       string `json:"name"`
	Constraint string `json:"constraint"`
}

// MetricValue is the SLA representation of a metric value.
type MetricValue struct {
	Key       string      `json:"key"`
	Action    string      `json:"action"`
	Namespace string      `json:"namespace"`
	Value     interface{} `json:"value"`
	DateTime  time.Time   `json:"datetime"`
}

func (v MetricValue) String() string {
	return fmt.Sprintf("{Key: %s, Value: %v, DateTime: %v}", v.Key, v.Value, v.DateTime)
}

// Violation is generated when a guarantee term is not fulfilled
type Violation struct {
	Id          string        `json:"id" bson:"_id"`
	AgreementId string        `json:"agreement_id"`
	Guarantee   string        `json:"guarantee"`
	Datetime    time.Time     `json:"datetime"`
	Constraint  string        `json:"constraint"`
	Values      []MetricValue `json:"values"`
	AppId       string        `json:"appID,omitempty"`
	Description string        `json:"description,omitempty"`
}

// SLAs is the type of an slice of SLA
type SLAs []SLA

// Violations is the type of an slice of Violation
type Violations []Violation

///////////////////////////////////////////////////////////////////////////////

// GetId returns the id of an agreement
func (a *SLA) GetId() string {
	return a.Id
}

// IsStarted is true if the agreement state is STARTED
func (a *SLA) IsStarted() bool {
	return a.State == STARTED
}

// IsTerminated is true if the agreement state is TERMINATED
func (a *SLA) IsTerminated() bool {
	return a.State == TERMINATED
}

// IsStopped is true if the agreement state is STOPPED
func (a *SLA) IsStopped() bool {
	return a.State == STOPPED
}

// IsValidTransition returns if the transition to newState is valid
func (a *SLA) IsValidTransition(newState State) bool {
	return a.State != TERMINATED
}

// Validate validates the consistency of an Agreement.
func (a *SLA) Validate(val Validator, mode ValidationMode) []error {
	return val.ValidateSLA(a, mode)
}

// Validate validates the consistency of an Assessment entity
func (as *Assessment) Validate(val Validator, mode ValidationMode) []error {
	return val.ValidateAssessment(as, mode)
}

// SetGuarantee is a helper function to set the assessment info of a guarantee term
func (as *Assessment) SetGuarantee(name string, value AssessmentGuarantee) {
	if as.Guarantees == nil {
		as.Guarantees = make(map[string]AssessmentGuarantee)
	}
	as.Guarantees[name] = value
}

// GetGuarantee is a helper to return the assessment info of a guarantee term.
//
// If empty, it returns a zero AssessmentGuarantee
func (as *Assessment) GetGuarantee(name string) AssessmentGuarantee {
	zero := AssessmentGuarantee{
		LastValues: LastValues{},
	}
	if as.Guarantees == nil {
		return zero
	}
	if _, ok := as.Guarantees[name]; !ok {
		return zero
	}
	return as.Guarantees[name]
}

// Validate validates the consistency of a Details entity
func (t *Details) Validate(val Validator, mode ValidationMode) []error {
	return val.ValidateDetails(t, mode)
}

// GetVariable returns the variable with name "varname".
//
// If not found, it returns a default value for the variable
// (i.e., Name and Metric equal to varname).
func (t *Details) GetVariable(varname string) (result Variable, ok bool) {
	for _, val := range t.Variables {
		if varname == val.Name {
			return val, true
		}
	}
	return Variable{Name: varname, Metric: varname}, false
}

// Validate validates the consistency of a Guarantee entity
func (g *Guarantee) Validate(val Validator, mode ValidationMode) []error {
	return val.ValidateGuarantee(g, mode)
}

// GetId returns the Id of a violation
func (v *Violation) GetId() string {
	return v.Id
}

// Validate validates the consistency of a Violation entity
func (v *Violation) Validate(val Validator, mode ValidationMode) []error {
	return val.ValidateViolation(v, mode)
}

// Normalize returns an always valid state: any different value from contained in States is STOPPED.
func (s State) Normalize() State {
	return normalizeState(s)
}
