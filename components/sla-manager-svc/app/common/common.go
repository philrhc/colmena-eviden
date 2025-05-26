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

package common

import (
	"os"
	"strconv"
)

/*
GetEnv get string environment variable value / return default if not found
*/
func GetEnv(name string, defaultval string) string {
	v := os.Getenv(name)
	if len(v) == 0 {
		return defaultval
	}
	return v
}

/*
GetEnv get int environment variable value / return default if not found
*/
func GetIntEnv(name string, defaultval int) int {
	v := os.Getenv(name)
	if len(v) == 0 {
		return defaultval
	}

	num, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return defaultval
	}
	return int(num)
}
