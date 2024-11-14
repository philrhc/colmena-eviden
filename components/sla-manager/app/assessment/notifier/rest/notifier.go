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

/* Implements notifier.NotifyViolations */
func (not _notifier) NotifyViolations(qos *model.SLA, result *amodel.Result) {

	vs := result.GetViolations()
	if len(vs) == 0 {
		return
	}

	info := violationInfo{
		Type:        "violation",
		AgreementID: qos.Id,
		Violations:  vs,
	}

	out, err1 := json.Marshal(info)
	if err1 == nil {
		logs.GetLogger().Infof("VIOLATOIN: " + string(out))
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
