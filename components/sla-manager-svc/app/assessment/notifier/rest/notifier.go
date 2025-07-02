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
package rest

import (
	amodel "colmena/sla-management-svc/app/assessment/model"
	"colmena/sla-management-svc/app/assessment/notifier"
	"colmena/sla-management-svc/app/model"

	"bytes"
	"encoding/json"
	"net/http"

	"colmena/sla-management-svc/app/common/cfg"
	"colmena/sla-management-svc/app/common/logs"

	"github.com/spf13/viper"
)

// path used in logs
const pathLOG string = "SLA > Assessment > Notifier > REST > "

type _notifier struct {
	url string
}

type violationInfo struct {
	Type          string            `json:"type"`
	AgreementID   string            `json:"agremeent_id"`
	GuaranteeName string            `json:"guarantee_name"`
	Violations    []model.Violation `json:"violations"`
}

// New constructs a REST Notifier
func New(config *viper.Viper) notifier.ViolationNotifier {

	logConfig(config)
	return _new(config.GetString(cfg.NotificationURLPropertyName))
}

func _new(url string) notifier.ViolationNotifier {
	return _notifier{
		url: url,
	}
}

func logConfig(config *viper.Viper) {
	logs.GetLogger().Info(pathLOG + "RestNotifier configuration\n" +
		"\t-----------------------------------------------------------------\n" +
		"\tURL (target of REST notifications): " + config.GetString(cfg.NotificationURLPropertyName) + "\n" +
		"\t-----------------------------------------------------------------")

}

/* Implements notifier.NotifyAllViolations */
func (not _notifier) NotifyAllViolations(results []model.OutputSLA) {
	out, err1 := json.Marshal(results)
	if err1 == nil {
		logs.GetLogger().Infof("VIOLATIONs: " + string(out))
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(results)

	_, err := http.Post(not.url, "application/json; charset=utf-8", b)

	if err != nil {
		logs.GetLogger().Error(pathLOG + "RestNotifier error: " + err.Error())
	} else {
		logs.GetLogger().Infof(pathLOG+"RestNotifier. Sent violations: %v", results)
	}
}

/* Implements notifier.NotifyAllStatuses */
func (not _notifier) NotifyAllStatuses(results []model.OutputSLA) {
	logs.GetLogger().Warn(pathLOG + "Function not implemented")
}

/* Implements notifier.NotifyViolations */
func (not _notifier) NotifyViolations(qos *model.SLA, result *amodel.Result) {

	vs := result.GetViolations()
	if len(vs) == 0 {
		return
	}

	/*
		Output model (SLA VIOLATION) example:

		{
			"serviceId": "test_service_id",
			"KPIs": [
				{
					"roleId": "",
					"query": "avg(processing_time)<5",
					"value": 10,
					"threshold": 5,
					"level": "",
					"violations": [],
					"total_violations": 1
				}
			]
		}
	*/
	info := model.OutputSLA{
		ServiceId: qos.Name,
		SLAId:     qos.Id,
		Kpis: []model.OutputSLAKpi{
			{
				RoleId:          qos.Details.Guarantees[0].Name,
				Query:           qos.Details.Guarantees[0].Constraint,
				Value:           vs[0].Values[0].Value,
				Level:           qos.Assessment.Level,
				Threshold:       qos.Assessment.Threshold, //qos.Details.Guarantees[0].Query,
				Violations:      vs,
				TotalViolations: qos.Assessment.TotalViolations,
			},
		},
	}

	out, err1 := json.Marshal(info)
	if err1 == nil {
		logs.GetLogger().Infof("VIOLATION: " + string(out))
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(info)

	_, err := http.Post(not.url, "application/json; charset=utf-8", b)

	if err != nil {
		logs.GetLogger().Error(pathLOG + "RestNotifier error: " + err.Error())
	} else {
		logs.GetLogger().Infof(pathLOG+"RestNotifier. Sent violations: %v", info)
	}
}

/* Implements notifier.NotifyStatus */
func (not _notifier) NotifyStatus(qos *model.SLA, result *amodel.Result) {
	var res interface{} = result.LastValues
	//logs.GetLogger().Debugf(pathLOG+" Value (1): ", res)

	if len(result.LastValues) > 0 {
		for key := range result.LastValues {
			//logs.GetLogger().Debugf(pathLOG+" key: ", key)
			if len(result.LastValues[key]) > 0 {
				for key2 := range result.LastValues[key] {
					//logs.GetLogger().Debugf(pathLOG+" key2: ", key2)
					if len(key2) > 0 {
						r, ok := result.LastValues[key][key2].Value.(float64)
						if !ok {
							logs.GetLogger().Error(pathLOG + " Value is not a number")
						} else {
							res = r
						}
						//logs.GetLogger().Debugf(pathLOG + " break")
						break
					}
				}
			}
		}
	}
	//logs.GetLogger().Debugf(pathLOG+"Value (2): ", res)

	info := model.OutputSLA{
		ServiceId: qos.Name,
		SLAId:     qos.Id,
		Kpis: []model.OutputSLAKpi{
			{
				RoleId:          qos.Details.Guarantees[0].Name,
				Query:           qos.Details.Guarantees[0].Constraint,
				Value:           res, //result.LastValues,
				Level:           qos.Assessment.Level,
				Threshold:       qos.Assessment.Threshold, //qos.Details.Guarantees[0].Query,
				TotalViolations: qos.Assessment.TotalViolations,
			},
		},
	}

	out, err1 := json.Marshal(info)
	if err1 == nil {
		logs.GetLogger().Infof("STATUS NOTIFICATION: " + string(out))
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(info)

	_, err := http.Post(not.url, "application/json; charset=utf-8", b)

	if err != nil {
		logs.GetLogger().Error(pathLOG + "RestNotifier error: " + err.Error())
	} else {
		logs.GetLogger().Infof(pathLOG+"RestNotifier. Sent status notification: %v", info)
	}
}
