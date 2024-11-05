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

// Package assessment contains the core code that evaluates the agreements.
// funcs.go contains the help functions used in evaluator file
package assessment

import (
	"fmt"
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
	logs.GetLogger().Debug(pathLOG + "[evaluateExpression] Evaluating expression ...")

	evalues := make(map[string]interface{})
	for key, value := range values {
		logs.GetLogger().Debug(pathLOG+"[evaluateExpression] key: "+key+", value: ", value)
		evalues[key] = value.Value
	}
	logs.GetLogger().Debug(pathLOG+"[evaluateExpression] (key, value): 	", evalues)

	result, err := expression.Evaluate(evalues)
	if err != nil {
		logs.GetLogger().Error(pathLOG + "[evaluateExpression] Error during the evaluation of the expression: " + err.Error())
		return nil, err
	}

	logs.GetLogger().Debug(pathLOG + "[evaluateExpression] Evaluating expression Vexpression=Vresult with values Vs")
	logs.GetLogger().Debug(pathLOG+"[evaluateExpression] Vexpression: 	", expression)
	logs.GetLogger().Debug(pathLOG+"[evaluateExpression] Vresult: 		", result)
	logs.GetLogger().Debug(pathLOG+"[evaluateExpression] Vs: ", values)

	if _, ok := result.(bool); ok {
		if !result.(bool) {
			return values, nil
		}
	} else {
		logs.GetLogger().Warn("[evaluateExpression] 'result' (from evaluation operation) is not a bool object.")

		str := fmt.Sprintf("%v", result)
		if strings.Contains(strings.ToLower(str), "false") {
			logs.GetLogger().Debug("[evaluateExpression] 'result' contains a 'false' value.")
			return values, nil
		} else if strings.Contains(strings.ToLower(str), "true") {
			logs.GetLogger().Debug("[evaluateExpression] 'result' contains a 'true' value.")
		} else {
			logs.GetLogger().Error(pathLOG + "[evaluateExpression] 'result' (from evaluation operation) is not a bool object.")
		}
	}

	return nil, err
}

/*
// Check guarantee "Warning" interfal if defined
func checkWarningInterval(agreement model.SLA, guarantee model.Guarantee, cfg Config, lastvalue amodel.ExpressionData) {
	not := cfg.Notifier
	for _, j := range guarantee.Importance {
		if j.Name != "Warning" {
			continue
		}
		var value float64
		for _, v := range lastvalue {
			value = v.Value.(float64)
		}

		expression, err := govaluate.NewEvaluableExpression(fmt.Sprintf("%f", value) + j.Constraint)
		if err != nil {
			log.Error(pathLOG + "funcs > [checkWarningInterval] Error while creating expression")
			continue
		}

		result, err := expression.Evaluate(nil)
		if err != nil {
			log.Error(pathLOG + "funcs > [checkWarningInterval] Error while evaluating expression")
			continue
		}

		if _, ok := result.(bool); ok {
			if result.(bool) {
				var result2 amodel.Result
				var evalgtresult amodel.EvaluationGtResult
				var violation = model.Violation{
					//Id:             agreement.Id,
					AgreementId:    agreement.Id,
					Guarantee:      guarantee.Name,
					Datetime:       cfg.Now,
					Constraint:     guarantee.Constraint,
					Values:         []model.MetricValue{},
					ImportanceName: "Warning",
					Importance:     -1,
					//AppId:          agreement.Details.Service,
					Description: "",
				}
				evalgtresult.Violations = append(evalgtresult.Violations, violation)
				result2.Violated = make(map[string]amodel.EvaluationGtResult)
				result2.Violated[guarantee.Name] = evalgtresult
				if not != nil {
					not.NotifyViolations(&agreement, &result2)
				}
			}
		} else {
			log.Warn("funcs > [checkWarningInterval] 'result' (from evaluation operation) is not a bool object.")

			str := fmt.Sprintf("%v", result)
			if strings.Contains(strings.ToLower(str), "false") {
				log.Debug("funcs > [checkWarningInterval] 'result' contains a 'false' value.")
			} else if strings.Contains(strings.ToLower(str), "true") {
				log.Debug("funcs > [checkWarningInterval] 'result' contains a 'true' value.")

				var result2 amodel.Result
				var evalgtresult amodel.EvaluationGtResult
				var violation = model.Violation{
					//Id:             agreement.Id,
					AgreementId:    agreement.Id,
					Guarantee:      guarantee.Name,
					Datetime:       cfg.Now,
					Constraint:     guarantee.Constraint,
					Values:         []model.MetricValue{},
					ImportanceName: "Warning",
					Importance:     -1,
					//AppId:          agreement.Details.Service,
					Description: "",
				}
				evalgtresult.Violations = append(evalgtresult.Violations, violation)
				result2.Violated = make(map[string]amodel.EvaluationGtResult)
				result2.Violated[guarantee.Name] = evalgtresult
				if not != nil {
					not.NotifyViolations(&agreement, &result2)
				}
			} else {
				log.Error(pathLOG + "funcs > [checkWarningInterval] 'result' (from evaluation operation) is not a bool object.")
			}
		}
	}
}
*/

// inTransientTime returns if the new violation detected occurs in the transient time
// of the guarantee term; i.e. last + transient < newviolation
func inTransientTime(newViolation time.Time, last *model.Violation, transientTime time.Duration) bool {
	// first violations are always considered out of the transient time
	if last == nil {
		return false
	}
	return newViolation.Before(last.Datetime.Add(transientTime))
}

// Prometheus interface for collect internal metrics from a Prometheus job: http://host:port/metrics
func updateInternalMetrics(v model.Violation) {
	//var intervalValue float64
	var intervalNames []string = []string{"Mild", "Serious", "Severe", "Catastrophic", "No violation"}
	for k := range intervalNames {
		//intervalValue = 0
		if v.ImportanceName == intervalNames[k] {
			/*intervalValue = v.Values[0].Value.(float64)
			metrics.CountViolation(intervalValue,
				map[string]string{
					"application": v.AppId,
					"agreement":   v.AgreementId,
					"metric":      v.Guarantee,
					//"importance":  v.IntervalName,
					"importance": intervalNames[k],
				},
			)*/
		}
		//metrics.AddSample(v.Values[0].Value.(float64),
		/*metrics.AddSample(intervalValue,
			map[string]string{
				"application": v.AppId,
				"agreement":   v.AgreementId,
				"metric":      v.Guarantee,
				//"importance":  v.IntervalName,
				"importance": intervalNames[k],
			},
		)*/
	}
}

/*
checkViolationLevel checks and sets violation level
*/
/*
func checkViolationLevel(a *model.SLA, gtv *[]model.Violation, v *model.Violation) {
	// violation leveles - Importance field
	interval := -1
	violationValue := 0.00

	ps := v // Pointer to the violation struct type

	// iterate "Importance" to get violation "level"
	for _, i := range a.Details.Guarantees {
		interval = -1
		if i.Name != v.Guarantee {
			continue
		}
		for _, j := range i.Importance {
			log.Trace(pathLOG + "vlevels [checkViolationLevel] ps.Values[0].Value.(float64) = " + fmt.Sprintf("%f", ps.Values[0].Value.(float64)))
			log.Trace(pathLOG + "vlevels [checkViolationLevel] j.Constraint = " + j.Constraint)
			expression, err := govaluate.NewEvaluableExpression(fmt.Sprintf("%f", ps.Values[0].Value.(float64)) + j.Constraint)
			if err != nil {
				log.Error(pathLOG + "vlevels [checkViolationLevel] Error while creating expression: " + err.Error())
				continue
			}
			result, err := expression.Evaluate(nil)
			if err != nil {
				log.Error(pathLOG + "vlevels [checkViolationLevel] Error while evaluating expression: " + err.Error())
				continue
			}

			str := fmt.Sprintf("%v", result)
			log.Debug(pathLOG + "vlevels [checkViolationLevel] " + fmt.Sprintf("%+v", j) + " " +
				"[" + fmt.Sprintf("%f", ps.Values[0].Value.(float64)) + " " + j.Constraint + "]=[" + str + "]")

			if _, ok := result.(bool); ok {
				if result.(bool) {
					violationValue = ps.Values[0].Value.(float64)
					interval++
				}
			} else {
				log.Warn(pathLOG + "vlevels [checkViolationLevel] 'result' (from evaluation operation) is not a bool object.")

				if strings.Contains(strings.ToLower(str), "false") {
					log.Debug(pathLOG + "vlevels [checkViolationLevel] 'result' contains a 'false' value.")
				} else if strings.Contains(strings.ToLower(str), "true") {
					log.Debug(pathLOG + "vlevels [checkViolationLevel] 'result' contains a 'true' value.")
					violationValue = ps.Values[0].Value.(float64)
					interval++
				} else {
					log.Error(pathLOG + "vlevels [checkViolationLevel] 'result' (from evaluation operation) is not a bool object.")
				}
			}

			switch interval {
			case 0:
				ps.ImportanceName = "Warning"
				ps.Importance = interval
			case 1:
				ps.ImportanceName = "Mild"
				ps.Importance = interval
			case 2:
				ps.ImportanceName = "Serious"
				ps.Importance = interval
			case 3:
				ps.ImportanceName = "Critical"
				ps.Importance = interval
			case 4:
				ps.ImportanceName = "Catastrophic"
				ps.Importance = interval
			default:
				ps.ImportanceName = "Default"
				ps.Importance = interval
			}
		}

		// add violation to list
		log.Debug(pathLOG + "vlevels [checkViolationLevel] Violation Value: [" + fmt.Sprintf("%+v", violationValue) + "]; ViolationType: [" + v.ImportanceName + "]")
		*gtv = append(*gtv, *v)

		// update / format content
		updateInternalMetrics(*v)
	}
}
*/
