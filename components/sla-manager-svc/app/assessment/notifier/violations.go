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
package notifier

import (
	assessment_model "colmena/sla-management-svc/app/assessment/model"
	"colmena/sla-management-svc/app/model"
)

type ViolationNotifier interface {
	NotifyViolations(agreement *model.SLA, result *assessment_model.Result)

	NotifyAllViolations(results []model.ColmenaOutputSLA)

	NotifyStatus(agreement *model.SLA, result *assessment_model.Result)

	NotifyAllStatuses(results []model.OutputSLA)
}
