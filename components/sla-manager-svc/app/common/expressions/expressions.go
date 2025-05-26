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
package expressions

import (
	"colmena/sla-management-svc/app/common/logs"
	"errors"
	"fmt"
	"strings"
)

// path used in logs
const pathLOG string = "SLA > Common "

const LABEL_MARK = "#LABELS#"

/*
Returns an array containg the expression '<metrics_query> <operator> <value>'
arr[0] = <metrics_query>
arr[1] = <operator>
arr[2] = <value>
*/
func getContraintParts(expr string) ([]string, error) {
	operators_arr := [6]string{"==", "<=", ">=", "!=", "<", ">"}
	for _, op := range operators_arr {
		if pos := strings.Index(expr, op); pos != -1 {
			arr := strings.Split(expr, string(op))
			if len(arr) != 2 {
				return nil, errors.New("not valid expression. Expression format: \"<metrics_query> <operator> <value>\"")
			} else {
				var res = make([]string, 3)
				res[0] = strings.Trim(arr[0], " ")
				res[1] = string(op)
				res[2] = strings.Trim(arr[1], " ")
				return res, nil
			}
		}
	}
	return nil, errors.New("no operator found. Valid operators: '==', '<=', '>=', '!=', '<', '>'. Expression format: \"<metrics_query> <operator> <value>\"")
}

/*
Insert mark for labels before time interval
*/
func setLabelsMark(metrics_query string) (string, error) {
	fmt.Println(metrics_query)
	posBracketIni := strings.Index(metrics_query, "[")

	totalInitBrackets := countOccurrences(metrics_query, '[')
	if posBracketIni > 0 && totalInitBrackets == 1 {
		// "avg_over_time(processing_time[5s])"
		logs.GetLogger().Debug(pathLOG + "==> periodo de tiempo")
		return "[" + strings.Replace(metrics_query, "[", LABEL_MARK+"[", 1) + "]", nil
	} else if posBracketIni == 0 && totalInitBrackets == 1 {
		// "[go_memstats_frees_total]"
		logs.GetLogger().Debug(pathLOG + "==> formato ok, sin periodo de tiempo")
		return strings.Replace(metrics_query, "]", LABEL_MARK+"]", 1), nil
	} else if posBracketIni == 0 && totalInitBrackets == 2 && strings.HasSuffix(metrics_query, "]") {
		// "[avg_over_time(processing_time[5s])]"
		logs.GetLogger().Debug(pathLOG + "==> formato ok. Buscar periodod de tiempo")
		return "[" + strings.Replace(metrics_query[1:], "[", LABEL_MARK+"[", 1), nil
	} else if posBracketIni < 0 {
		// No '[', ']'
		// buscamos ')'
		posParenthEnd := strings.Index(metrics_query, ")")
		if posParenthEnd > 0 {
			logs.GetLogger().Debug(pathLOG + "==> formato ok, sin brackets")
			return "[" + strings.Replace(metrics_query, ")", LABEL_MARK+")", 1), nil
		} else {
			logs.GetLogger().Debug(pathLOG + "==> formato ok, sin brackets")
			return "[" + metrics_query + LABEL_MARK + "]", nil
		}
	}

	return "", errors.New("not valid 'metrics_query' expression. Expression format: '[<expression>]' or '<expression>[interval]' or <expression> ")
}

/*
get total ocurrences of a character in a string
*/
func countOccurrences(s string, char rune) int {
	count := 0
	for _, c := range s {
		if c == char {
			count++
		}
	}
	return count
}

/*
Objective:

	from "avg_over_time(processing_time[5s]) < 1"
	to "[avg_over_time(processing_time#LABELS#[5s])] < 1"
	where #LABELS# will be repalced by {building='BSC'}
*/
func CheckAndParseConstraint(constraint string) (string, error) {
	logs.GetLogger().Debug(pathLOG + "[CheckAndParseConstraint] Checking expression: " + constraint)
	res, err := getContraintParts(constraint)
	if err != nil {
		logs.GetLogger().Error(pathLOG+"[CheckAndParseConstraint] ", err)
	} else {
		logs.GetLogger().Debug(pathLOG + "[CheckAndParseConstraint] " + strings.Join(res, " // "))
		res2, err := setLabelsMark(res[0])
		if err != nil {
			logs.GetLogger().Error(pathLOG+"[CheckAndParseConstraint] ", err)
		} else {
			logs.GetLogger().Debug(pathLOG + "[CheckAndParseConstraint] " + res2)
			return res2 + " " + res[1] + " " + res[2], nil
		}
	}

	return "", err
}

/*
 */
func ParseConstraint(constraint string) (string, error) {
	// example: "[avg_over_time(go_goroutines[60m])] < 50000"
	logs.GetLogger().Debug(pathLOG + "[ParseConstraint] Checking and parsing constraint expression " + constraint + " ...")
	logs.GetLogger().Debug(pathLOG + "[ParseConstraint] Constraint expression format: [expression] '<'/'='/'>' value")

	expr := strings.TrimSpace(constraint)
	// example: "[avg_over_time(go_goroutines[60m])] < 50000"
	logs.GetLogger().Debug(pathLOG + "[ParseConstraint] expr: " + expr)

	posBracketIni := strings.Index(expr, "[")
	if posBracketIni != 0 {
		return "", errors.New("bad constraint expression: No initial bracket found on left side. \n" +
			"Expected Constraint expression format: '['<expression>']' '<'/'='/'>' <value> \n" +
			"Example: [avg_over_time(go_goroutines[60m])] < 50000") // ERROR
	} else {
		exprsubstring1 := expr[1:]
		// example: "avg_over_time(go_goroutines[60m])] < 50000"
		logs.GetLogger().Debug(pathLOG + "[ParseConstraint] exprsubstring1: " + exprsubstring1)

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
				logs.GetLogger().Debug(pathLOG + "[ParseConstraint] exprsubstring2: " + exprsubstring2)

				exprsubstring3 := exprsubstring1[posBracketEnd+1:]
				logs.GetLogger().Debug(pathLOG + "[ParseConstraint] exprsubstring3: " + exprsubstring3)

				// https://www.w3schools.com/tags/ref_urlencode.asp
				exprsubstring2 = strings.ReplaceAll(exprsubstring2, "[", "%5B")
				exprsubstring2 = strings.ReplaceAll(exprsubstring2, "]", "%5D")

				exprStrFinal := "[" + exprsubstring2 + "]" + exprsubstring3
				logs.GetLogger().Debug(pathLOG + "[ParseConstraint] exprStrFinal: " + exprStrFinal)

				return exprStrFinal, nil
			}
		}
	}
}
