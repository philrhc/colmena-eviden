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
	"colmena/sla-management-svc/app/common/logs"
	"errors"
	"strings"
)

func parseConstraint(constraint string) (string, error) {
	// example: "[avg_over_time(go_goroutines[60m])] < 50000"
	logs.GetLogger().Debug(pathLOG + "[parseConstraint] Checking and parsing constraint expression " + constraint + " ...")
	logs.GetLogger().Debug(pathLOG + "[parseConstraint] Constraint expression format: [expression] '<'/'='/'>' value")

	expr := strings.TrimSpace(constraint)
	// example: "[avg_over_time(go_goroutines[60m])] < 50000"
	logs.GetLogger().Debug(pathLOG + "[parseConstraint] expr: " + expr)

	posBracketIni := strings.Index(expr, "[")
	if posBracketIni != 0 {
		return "", errors.New("bad constraint expression: No initial bracket found on left side. \n" +
			"Expected Constraint expression format: '['<expression>']' '<'/'='/'>' <value> \n" +
			"Example: [avg_over_time(go_goroutines[60m])] < 50000") // ERROR
	} else {
		exprsubstring1 := expr[1:]
		// example: "avg_over_time(go_goroutines[60m])] < 50000"
		logs.GetLogger().Debug(pathLOG + "[parseConstraint] exprsubstring1: " + exprsubstring1)

		posBracketIni2 := strings.Index(exprsubstring1, "[")
		if posBracketIni2 < 0 {
			// example: "[avg_over_time(go_goroutines[60m])] < 50000"
			return constraint, nil
		} else {
			posBracketEnd := strings.LastIndex(exprsubstring1, "]")
			if posBracketEnd < 0 {
				return "", errors.New("bad constraint expression: No final bracket found. \n" +
					"Expected Constraint expression format: '['<expression>']' '<'/'='/'>' <value> \n" +
					"Example: [avg_over_time(go_goroutines[60m])] < 50000") // ERROR
			} else {

				exprsubstring2 := exprsubstring1[0:posBracketEnd]
				logs.GetLogger().Debug(pathLOG + "[parseConstraint] exprsubstring2: " + exprsubstring2)

				exprsubstring3 := exprsubstring1[posBracketEnd+1:]
				logs.GetLogger().Debug(pathLOG + "[parseConstraint] exprsubstring3: " + exprsubstring3)

				// https://www.w3schools.com/tags/ref_urlencode.asp
				exprsubstring2 = strings.ReplaceAll(exprsubstring2, "[", "%5B")
				exprsubstring2 = strings.ReplaceAll(exprsubstring2, "]", "%5D")

				exprStrFinal := "[" + exprsubstring2 + "]" + exprsubstring3
				logs.GetLogger().Debug(pathLOG + "[parseConstraint] exprStrFinal: " + exprStrFinal)

				return exprStrFinal, nil
			}
		}
	}
}
