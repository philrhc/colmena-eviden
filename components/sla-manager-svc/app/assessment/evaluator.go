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
	"time"

	amodel "colmena/sla-management-svc/app/assessment/model"
	"colmena/sla-management-svc/app/assessment/monitor"
	"colmena/sla-management-svc/app/assessment/notifier"
	"colmena/sla-management-svc/app/common/logs"
	"colmena/sla-management-svc/app/model"

	"github.com/Knetic/govaluate"
)

// path used in logs
const pathLOG string = "SLA > App > Assessment [EVALUATION PHASE] "

/*
Config contains the configuration for an assessment process (i.e. global config)
*/
type Config struct {
	// Now is the time considered the current time. In general terms, metrics are retrieved from the adapter from the last measure to `now`.
	// If the monitoring have some delay storing metrics, now could be shifted some minutes to the past.
	Now time.Time

	// Repo is the repository where to load/store entities
	Repo model.IRepository

	// Adapter is the monitoring adapter where to get metrics from
	Adapter monitor.MonitoringAdapter

	// Notifier receives the violations and notifies them (send by REST, store to DB...)
	Notifier notifier.ViolationNotifier

	// Transient is time to wait until a new violation of a GT can be raised again (default value is zero)
	Transient time.Duration
}

/*
AssessActiveQoSDefinitions will get the active QoS from the provided repository and assess them, notifying about violations with
the provided notifier. Mandatory fields filled in cfg are Repo, Adapter and Now.
*/
func AssessActiveQoSDefinitions(cfg Config) {
	repo := cfg.Repo
	not := cfg.Notifier

	// Retrieve all active QoS definitions
	qosdefs, err := repo.GetSLAsByState(model.STARTED, model.STOPPED)

	if err != nil {
		logs.GetLogger().Error(pathLOG+"[AssessActiveQoSDefinitions] Error getting active qos definitions: %s", err.Error())
	} else {
		logs.GetLogger().Infof(pathLOG+"[AssessActiveQoSDefinitions] [%d SLAs definitions to evaluate]", len(qosdefs))

		grouped_qosdefs, err := groupSLAsByServiceId(qosdefs, repo)
		if err == nil {

			var violations []model.OutputSLA // list of all violations
			var statuses []model.OutputSLA   // list of all violations

			// iterate SLA evaluation results
			for _, qosdefs2 := range grouped_qosdefs {
				if len(qosdefs2) > 0 {
					logs.GetLogger().Infof(pathLOG + "[AssessActiveQoSDefinitions] Evaluating service [" + qosdefs2[0].Name + "]")

					for _, qosd := range qosdefs2 {
						if qosd.State != model.STARTED {
							logs.GetLogger().Error(pathLOG + "[AssessActiveQoSDefinitions] SLA not started")
						} else {
							// do QoS assessment
							logs.GetLogger().Debug(pathLOG+"[AssessActiveQoSDefinitions] SLA Assessment ", qosd.Id)

							result, totalResults := AssessQoS(&qosd, cfg)
							qosd.Assessment.TotalExecutions += 1

							// violation?
							violation := not != nil && len(result.Violated) > 0

							if violation {
								qosd.Assessment.TotalViolations += 1
								qosd.Assessment.Violated = true

								violation_result := GenerateViolationOutput(qosd, result)
								if violation_result.ServiceId != "" {
									violations = append(violations, violation_result)
								}

							} else {
								qosd.Assessment.Violated = false
							}

							// check and set violation levels
							checkViolationLevel(&qosd, totalResults)

							// notify violations or status
							if violation {
								//not.NotifyViolations(&qosd, &result)
							} else {
								not.NotifyStatus(&qosd, &result)
							}

							// update SLA
							repo.UpdateSLA(&qosd)
						}
					}
				} else {
					logs.GetLogger().Error(pathLOG + "[AssessActiveQoSDefinitions] Error: emty service SLA found!")
				}
			}

			if len(violations) > 0 {
				not.NotifyAllViolations(violations)
			}

			if len(statuses) > 0 {
				// TODO send all notifications

			}
		} else {
			logs.GetLogger().Error(pathLOG+"[AssessActiveQoSDefinitions] Error getting SLAs by id: %s", err.Error())
		}
	}
}

/*
AssessQoS is the process that assess a QoS definition. The process is:
 1. Check expiration date
 2. Evaluate metrics defined in Guarantees if QoS is started
 3. Set LastExecution time.

The output is:
  - parameter a is modified
  - evaluation results are the function return (violated metrics and raised violations).
  - a guarantee term is filled in the result only if there are violations.

The function results are not persisted. The output must be persisted/handled accordingly.
E.g.: QoSDefinition and Violations must be persisted to DB. Violations must be notified to observers
*/
func AssessQoS(a *model.SLA, cfg Config) (amodel.Result, int) {
	now := cfg.Now

	if a.Expiration != nil && a.Expiration.Before(now) {
		a.State = model.TERMINATED // QoSDefinition has expired
		logs.GetLogger().Debug(pathLOG + "[AssessQoS] QoS with ID [" + a.Id + "] has EXPIRED")
	}

	totalResults := 0
	if a.State == model.STARTED {
		logs.GetLogger().Debug(pathLOG+"[AssessQoS] Assessing QoS with ID: ", a.Id)

		// evaluates the guarantee terms defined in the QoS definition
		result, t, err := EvaluateGuaranteeTerms(a, cfg)
		if err != nil {
			logs.GetLogger().Warn(pathLOG + "[AssessQoS] Error evaluating QoSDefinition " + a.Id + ": " + err.Error())
			return amodel.Result{}, 0
		}
		updateAssessment(a, result, now) // updates QoSDefinition with last results
		totalResults = t

		if totalResults > 0 {
			logs.GetLogger().Debug(pathLOG+"[AssessQoS] QoS with ID ["+a.Id+"] has VIOLATIONS: ", result)
		}

		return result, totalResults
	}

	return amodel.Result{}, totalResults
}

/*
EvaluateGuaranteeTerms evaluates the guarantee terms of a QoS definition. The metric values are retrieved from a MonitoringAdapter.
The MonitoringAdapter must feed the process correctly (e.g. if the constraint of a guarantee term is of the type "A>B && C>D", the
MonitoringAdapter must supply pairs of values).
*/
func EvaluateGuaranteeTerms(a *model.SLA, cfg Config) (amodel.Result, int, error) {
	ma := cfg.Adapter.Initialize(a)
	now := cfg.Now

	result := amodel.Result{
		Violated:      map[string]amodel.EvaluationGtResult{},
		LastValues:    map[string]amodel.ExpressionData{},
		LastExecution: map[string]time.Time{},
	}
	gts := a.Details.Guarantees

	totalResults := 0

	for _, gt := range gts {
		// evaluates a guarantee term of the QoS Definition
		failed, lastvalues, _, err := EvaluateGuarantee(a, gt, ma, cfg)
		if err != nil {
			logs.GetLogger().Warn(pathLOG + "[EvaluateGuaranteeTerms] Error evaluating expression " + gt.Constraint + ": " + err.Error())
			return amodel.Result{}, 0, err
		}

		if len(failed) > 0 {
			// VIOLATIONS
			violations := EvaluateGtViolations(a, gt, failed, cfg.Transient) // Evaluates violation
			gtResult := amodel.EvaluationGtResult{
				Metrics:    failed,
				Violations: violations,
			}
			result.Violated[gt.Name] = gtResult
			totalResults = len(violations)
		}
		result.LastValues[gt.Name] = lastvalues
		result.LastExecution[gt.Name] = now

		logs.GetLogger().Debug(pathLOG+"[EvaluateGuaranteeTerms] result: ", result)

	}
	return result, totalResults, nil
}

/*
EvaluateGuarantee evaluates a guarantee term of a QoS Definition (see EvaluateGuaranteeTerms) and returns the metrics that failed the GT constraint.
*/
func EvaluateGuarantee(a *model.SLA, gt model.Guarantee, ma monitor.MonitoringAdapter,
	cfg Config) (failed []amodel.ExpressionData, last amodel.ExpressionData, totalResults int, err error) {

	logs.GetLogger().Debug(pathLOG + "[EvaluateGuarantee] Evaluating Guarantee [" + gt.Name + "] of QoS with ID [" + a.Id + "]; Expression: " + gt.Constraint)
	totalResults = 0
	failed = make(amodel.GuaranteeData, 0, 1)

	constraintParsedExpr, err := parseConstraint(gt.Constraint)
	if err != nil {
		logs.GetLogger().Error(pathLOG+"[EvaluateGuarantee] Error parsing expression: ", gt.Constraint)
		return nil, nil, totalResults, err
	}

	expression, err := govaluate.NewEvaluableExpression(constraintParsedExpr) //constraintParsedExpr) //gt.Constraint)
	if err != nil {
		logs.GetLogger().Error(pathLOG+"[EvaluateGuarantee] Error parsing expression: ", constraintParsedExpr)
		return nil, nil, totalResults, err
	}

	logs.GetLogger().Debug(pathLOG + "[EvaluateGuarantee] Getting values from monitor ...")
	values := ma.GetValues(gt, expression.Vars(), cfg.Now)

	if len(values) == 0 {
		logs.GetLogger().Warn(pathLOG+"[EvaluateGuarantee] No values found for Guarantee ["+gt.Name+"] of agreement with ID: ", a.Id)
	} else {
		logs.GetLogger().Debug(pathLOG+"[EvaluateGuarantee] Total values returned from Monitor ["+a.Id+", "+gt.Name+"]: ", len(values))
	}

	for _, value := range values {
		aux, err := evaluateExpression(expression, value)
		if err != nil {
			logs.GetLogger().Warn("[EvaluateGuarantee] Error evaluating expression " + gt.Constraint + ": " + err.Error())
			return nil, nil, totalResults, err
		}
		if aux != nil {
			failed = append(failed, aux)
		}
	}
	if len(values) > 0 {
		last = values[len(values)-1]
	}
	totalResults = len(failed)

	return failed, last, totalResults, nil
}

/*
EvaluateGtViolations creates violations for the detected violated metrics in EvaluateGuarantee
*/
func EvaluateGtViolations(a *model.SLA, gt model.Guarantee, violated amodel.GuaranteeData, transientTime time.Duration) []model.Violation {
	gtv := make([]model.Violation, 0, len(violated))
	lastViolation := a.Assessment.GetGuarantee(gt.Name).LastViolation

	for _, tuple := range violated {
		// build values map and find newer metric
		var d *time.Time
		var values = make([]model.MetricValue, 0, len(tuple))
		for _, m := range tuple {
			values = append(values, m)
			if d == nil || m.DateTime.After(*d) {
				d = &m.DateTime
			}
		}
		if inTransientTime(*d, lastViolation, transientTime) {
			logs.GetLogger().Debug(pathLOG+"[EvaluateGtViolations] Skipping failed metrics %v; last=%s transient=%d newTime=%s", tuple, lastViolation, transientTime, *d)
			continue
		}

		// VIOLATION object
		// with default violation leveles - Importance fields: (intervalName := "Default"), (interval := -1)
		v := model.Violation{
			AgreementId: a.Id,
			Guarantee:   gt.Name,
			Datetime:    *d,
			Constraint:  gt.Constraint,
			Values:      values,
			AppId:       a.Id,
			Description: "",
		}

		lastViolation = &v // update last violation value

		// check violation level: e.g., Mild, Serious etc.
		//checkViolationLevel(a, &gtv, &v)

		gtv = append(gtv, v)
	}
	logs.GetLogger().Debug(pathLOG+"[EvaluateGtViolations] Violations list content: ", gtv)
	return gtv
}

/*
BuildRetrievalItems returns the RetrievalItems to be passed to an EarlyRetriever.
*/
func BuildRetrievalItems(a *model.SLA, gt model.Guarantee, varnames []string, to time.Time) []monitor.RetrievalItem {
	result := make([]monitor.RetrievalItem, 0, len(varnames))

	defaultFrom := getDefaultFrom(a, gt)
	for _, name := range varnames {
		v, _ := a.Details.GetVariable(name)
		from := getFromForVariable(v, defaultFrom, to)
		item := monitor.RetrievalItem{
			Guarantee: gt,
			Var:       v,
			From:      from,
			To:        to,
		}
		result = append(result, item)
	}
	return result
}
