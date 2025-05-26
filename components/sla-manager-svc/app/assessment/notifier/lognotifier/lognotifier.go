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

// Package lognotifier contains a simple ViolationsNotifier that just logs violations.
package lognotifier

import (
	assessment_model "colmena/sla-management-svc/app/assessment/model"
	"colmena/sla-management-svc/app/common/logs"
	"colmena/sla-management-svc/app/model"
	"fmt"
)

// path used in logs
const pathLOG string = "SLA > Assessment > Notifier > Logs > "

// LogNotifier logs violations
type LogNotifier struct {
}

type violationInfo struct {
	ServiceId string `json:"serviceId"`
	RoleId    string `json:"roleId"`
}

// NotifyViolations implements ViolationNotifier interface
func (n LogNotifier) NotifyViolations(qos *model.SLA, result *assessment_model.Result) {

	logs.GetLogger().Info(pathLOG + "Checking violations ...")

	for k, v := range result.Violated {

		if len(v.Violations) > 0 {
			logs.GetLogger().Info(pathLOG + "Violation of agreement: " + qos.Id + "; Failed guarantee: " + k)

			for _, vi := range v.Violations {
				notification := violationInfo{
					ServiceId: qos.Name,
					RoleId:    vi.Guarantee,
				}

				strNotification := fmt.Sprintf("%+v", notification)
				strGuarantee := fmt.Sprintf("%+v", vi)

				logs.GetLogger().Infof(pathLOG + "Failed guarantee: " + strGuarantee)
				logs.GetLogger().Infof(pathLOG + "Sent violations: " + strNotification)
			}
		}
	}
}

/* Implements notifier.NotifyStatus */
func (n LogNotifier) NotifyStatus(qos *model.SLA, result *assessment_model.Result) {
	logs.GetLogger().Warn(pathLOG + "Function not implemented")
}

/* Implements notifier.NotifyAllViolations */
func (n LogNotifier) NotifyAllViolations(results []model.OutputSLA) {
	logs.GetLogger().Warn(pathLOG + "Function not implemented")
}

/* Implements notifier.NotifyAllStatuses */
func (n LogNotifier) NotifyAllStatuses(results []model.OutputSLA) {
	logs.GetLogger().Warn(pathLOG + "Function not implemented")
}
