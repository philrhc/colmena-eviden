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
	"fmt"
	"slices"
	"strings"
	"time"

	amodel "colmena/sla-management-svc/app/assessment/model"
	"colmena/sla-management-svc/app/common/logs"
	"colmena/sla-management-svc/app/model"

	"github.com/Knetic/govaluate"
)

// updateAssessment
func updateAssessment(a *model.SLA, result amodel.Result, now time.Time) {
	if a.Assessment.FirstExecution.IsZero() {
		a.Assessment.FirstExecution = now
	}
	a.Assessment.LastExecution = now

	for _, gt := range a.Details.Guarantees {
		gtname := gt.Name
		last := result.LastValues[gtname]

		violations := []model.Violation{}
		if violated, ok := result.Violated[gtname]; ok {
			violations = violated.Violations
		}
		updateAssessmentGuarantee(a, gtname, last, violations, now)
	}
}

// updateAssessmentGuarantee
func updateAssessmentGuarantee(a *model.SLA, gtname string, last amodel.ExpressionData,
	violations []model.Violation, now time.Time) {

	ag := a.Assessment.GetGuarantee(gtname)
	ag.LastExecution = now
	if ag.FirstExecution.IsZero() {
		ag.FirstExecution = now
	}
	for _, v := range last {
		ag.LastValues[v.Key] = v
	}
	if len(violations) > 0 {
		ag.LastViolation = &violations[len(violations)-1]
	}
	a.Assessment.SetGuarantee(gtname, ag)
}

// getDefaultFrom
func getDefaultFrom(a *model.SLA, gt model.Guarantee) time.Time {
	var defaultFrom = a.Assessment.GetGuarantee(gt.Name).LastExecution
	if defaultFrom.IsZero() {
		defaultFrom = a.Assessment.LastExecution
	}
	if defaultFrom.IsZero() {
		defaultFrom = a.Creation
	}
	return defaultFrom
}

// getFromForVariable returns the interval start for the query to monitoring.
// If the variable is aggregated, it depends on the aggregation window.
// If not, returns defaultFrom (which should be the last time the guarantee term was evaluated)
func getFromForVariable(v model.Variable, defaultFrom, to time.Time) time.Time {
	if v.Aggregation != nil && v.Aggregation.Window != 0 {
		return to.Add(-time.Duration(v.Aggregation.Window) * time.Second)
	}
	return defaultFrom
}

/*
evaluateExpression evaluate a GT expression at a single point in time with a tuple of metric values (one value per variable in GT expresssion)
The result is: the values if the expression is false (i.e., the failing values), or nil if expression was true
*/
func evaluateExpression(expression *govaluate.EvaluableExpression, values amodel.ExpressionData) (amodel.ExpressionData, error) {
	logs.GetLogger().Debug(pathLOG+"[evaluateExpression] Evaluating expression: '", expression, "' with the following values: ", values)

	evalues := make(map[string]interface{})
	for key, value := range values {
		evalues[key] = value.Value
	}

	result, err := expression.Evaluate(evalues)
	if err != nil {
		logs.GetLogger().Error(pathLOG + "[evaluateExpression] Error during the evaluation of the expression: " + err.Error())
		return nil, err
	}
	logs.GetLogger().Debug(pathLOG+"[evaluateExpression] Expression evaluation result: ", result)

	if _, ok := result.(bool); ok {
		if !result.(bool) {
			logs.GetLogger().Debug("[evaluateExpression] Returning content: ", values)
			return values, nil
		}
	} else {
		logs.GetLogger().Warn("[evaluateExpression] 'result' (from evaluation operation) is not a bool object.")

		str := fmt.Sprintf("%v", result)
		if strings.Contains(strings.ToLower(str), "false") {
			logs.GetLogger().Warn("[evaluateExpression] 'result' contains a 'false' value.")
			return values, nil
		} else if strings.Contains(strings.ToLower(str), "true") {
			logs.GetLogger().Warn("[evaluateExpression] 'result' contains a 'true' value.")
		} else {
			logs.GetLogger().Error(pathLOG + "[evaluateExpression] 'result' (from evaluation operation) is not a bool object.")
		}
	}

	return nil, err
}

/*
1) Primera vez que una KPI se infringe -> Level = Broken
2) Después de X veces [seguidas] que se ha infringido -> Level = Critical
3) Primera vez que una KPI se cumple [después de estar Broken] -> Level = Met
4) Después de Y veces [seguidas] que se ha cumplido -> Level = Desired
5) Sí ha cambiado de KPI met a KPI broken Z veces -> Level = Unstable

Levels: Broken, Critical, Met, Desired, Unstable, Unknown
*/
func checkViolationLevel(qos *model.SLA, totalResults int) {
	logs.GetLogger().Debug(pathLOG+"[checkViolationLevel] totalResults: ", totalResults)
	if totalResults == 0 {
		qos.Assessment.Level = model.ASSESSMENT_LEVEL_UNKNOWN
		return
	}

	if qos.Assessment.Violated {
		qos.Assessment.XCounter += 1
		if qos.Assessment.YCounter > 0 {
			qos.Assessment.ZCounter += 1
			qos.Assessment.YCounter = 0
		}
	} else {
		qos.Assessment.YCounter += 1
	}

	if qos.Assessment.Violated && qos.Assessment.XCounter == 1 {
		// 1) Primera vez que una KPI se infringe -> Level = Broken
		qos.Assessment.Level = model.ASSESSMENT_LEVEL_BROKEN
	} else if qos.Assessment.Violated && qos.Assessment.XCounter >= model.DEFAULT_ASSESSMENT_X {
		// 2) Después de X veces [seguidas] que se ha infringido -> Level = Critical
		qos.Assessment.Level = model.ASSESSMENT_LEVEL_CRITICAL
	} else if !qos.Assessment.Violated && qos.Assessment.YCounter == 1 {
		// 3) Primera vez que una KPI se cumple [después de estar Broken] -> Level = Met
		qos.Assessment.Level = model.ASSESSMENT_LEVEL_MET
	} else if !qos.Assessment.Violated && qos.Assessment.YCounter >= model.DEFAULT_ASSESSMENT_Y {
		// 4) Después de Y veces [seguidas] que se ha cumplido -> Level = Desired
		qos.Assessment.Level = model.ASSESSMENT_LEVEL_DESIRED
	} else if qos.Assessment.Violated && qos.Assessment.ZCounter >= model.DEFAULT_ASSESSMENT_Z {
		// 5) Sí ha cambiado de KPI met a KPI broken Z veces -> Level = Unstable
		qos.Assessment.Level = model.ASSESSMENT_LEVEL_UNSTABLE
	} else {
		qos.Assessment.Level = model.ASSESSMENT_LEVEL_UNKNOWN
	}

}

// inTransientTime returns if the new violation detected occurs in the transient time
// of the guarantee term; i.e. last + transient < newviolation
func inTransientTime(newViolation time.Time, last *model.Violation, transientTime time.Duration) bool {
	// first violations are always considered out of the transient time
	if last == nil {
		return false
	}
	return newViolation.Before(last.Datetime.Add(transientTime))
}

// func groupSLAsByServiceId
func groupSLAsByServiceId(qosdefs model.SLAs, repo model.IRepository) ([]model.SLAs, error) {
	var ids []string
	var res []model.SLAs

	for _, qosd := range qosdefs {
		if len(ids) == 0 || !slices.Contains(ids, qosd.Name) {
			ids = append(ids, qosd.Name)
		}
	}

	for _, id := range ids {
		// Retrieve all active QoS definitions
		r, err := repo.GetSLAsByName(id)
		if err == nil {
			res = append(res, r)
		} else {
			return nil, err
		}
	}

	return res, nil
}
