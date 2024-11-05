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
package model

import (
	"encoding/json"
	"os"
)

// ReadAgreement returns the agreement read from the file pointed by path.
// The CWD is the location of the test.
//
// Ex:
//
//	a, err := readAgreement("testdata/a.json")
//	if err != nil {
//	  t.Errorf("Error reading agreement: %v", err)
//	}
func ReadAgreement(path string) (SLA, error) {
	res, err := readEntity(path, new(SLA))
	a := res.(*SLA)

	return *a, err
}

// readEntity
func readEntity(path string, result interface{}) (interface{}, error) {

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return result, err
	}
	json.NewDecoder(f).Decode(&result)
	return result, err
}
