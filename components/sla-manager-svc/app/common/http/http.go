/*
  COLMENA-DESCRIPTION-SERVICE
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
package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"colmena/sla-management-svc/app/common/logs"
)

// path used in logs
const pathLOG string = "QAA > App > Common > HTTP > "

///////////////////////////////////////////////////////////////////////////////

// creates the request's body' from a JSON
func httpJSONBody(bodyJSON interface{}) io.Reader {
	bodyBytes, err := json.Marshal(bodyJSON)
	if err != nil {
		logs.GetLogger().Error(pathLOG+"[httpJSONBody] ERROR (1)", err)
		return nil
	}
	return bytes.NewReader(bodyBytes)
}

// creates the request's body' from a string
func httpRawDataBody(bodyRawData string) io.Reader {
	return bytes.NewReader([]byte(bodyRawData))
}

// httpRequest prepares and executes the HTTP request //// bodyJSON interface{}
func httpRequest(httpMethod string, url string, auth bool, connToken string, body io.Reader) (int, []byte, error) {
	logs.GetLogger().Info(pathLOG + "[httpRequest] " + httpMethod + " request [" + url + "], auth [" + strconv.FormatBool(auth) + "] ...")

	// create request with headers and body
	req, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		logs.GetLogger().Error(pathLOG+"[httpRequest] ERROR (2)", err)
		return 0, nil, err
	}

	// Content-Type: json / json-patch+json
	if httpMethod == "PATCH" {
		logs.GetLogger().Info(pathLOG + "[httpRequest] Content-Type = application/json-patch+json")
		req.Header.Set("Content-Type", "application/json-patch+json")
	} else {
		logs.GetLogger().Info(pathLOG + "[httpRequest] Content-Type = application/json")
		req.Header.Set("Content-Type", "application/json")
	}

	// Authorization header
	if auth == true {
		logs.GetLogger().Info(pathLOG + "[httpRequest] Using Authorization Bearer ...")
		req.Header.Set("Authorization", "Bearer "+connToken)
	}

	// CLIENT
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// execute HTTP request
	resp, err := client.Do(req)
	if err != nil {
		logs.GetLogger().Error(pathLOG+"[httpRequest] ERROR (3)", err)
		return 0, nil, err
	}
	defer resp.Body.Close()

	// get data from response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.GetLogger().Error(pathLOG+"[httpRequest] ERROR (4)", err)
		return resp.StatusCode, nil, err
	} else if resp.StatusCode >= 400 { // check errors => StatusCode
		logs.GetLogger().Info(pathLOG + "[httpRequest] ERROR (5) StatusCode >= 400")
		return resp.StatusCode, nil, errors.New(pathLOG + "[httpRequest] HTTP STATUS: (" + strconv.Itoa(resp.StatusCode) + ") " + http.StatusText(resp.StatusCode) + "")
	}

	logs.GetLogger().Info(pathLOG + "[httpRequest] HTTP STATUS: (" + strconv.Itoa(resp.StatusCode) + ") " + http.StatusText(resp.StatusCode))

	return resp.StatusCode, data, nil
}

///////////////////////////////////////////////////////////////////////////////
// GET

/*
Get generic GET request
*/
func Get(url string, auth bool, connToken string) (int, []byte, error) {
	return httpRequest("GET", url, auth, connToken, nil)
}

/*
Get generic GET request with payload
*/
func GetWithPayload(url string, auth bool, connToken string, payload string) (int, []byte, error) {
	return httpRequest("GET", url, auth, connToken, strings.NewReader(payload))
}

/*
GetStruct GET request that returns a struct of type 'map[string]interface{}'
*/
func GetStruct(url string, auth bool, connToken string) (int, map[string]interface{}, error) {
	logs.GetLogger().Info(pathLOG + "[GetStruct] GET request [" + url + "] ...")

	status, data, err := Get(url, auth, connToken)
	if err != nil {
		logs.GetLogger().Error(pathLOG+"[GetStruct] ERROR (1)", err)
		return status, nil, err
	}

	// create json
	var objmap map[string]interface{}
	if err := json.Unmarshal(data, &objmap); err != nil {
		logs.GetLogger().Error(pathLOG+"[GetStruct] ERROR (2)", err)
		return status, nil, err
	}

	return status, objmap, nil
}

/*
GetString GET request that returns a string (response)
*/
func GetString(url string, auth bool, connToken string) (int, string, error) {
	logs.GetLogger().Info(pathLOG + "[GetString] GET request [" + url + "] ...")

	status, data, err := Get(url, auth, connToken)
	if err != nil {
		logs.GetLogger().Error(pathLOG+"[GetString] ERROR (1)", err)
		return status, "", err
	}

	return status, string(data), nil
}

///////////////////////////////////////////////////////////////////////////////
// POST

/*
PostRawData Generic POST request
*/
func PostRawData(url string, auth bool, connToken string, bodyRawData string) (int, []byte, error) {
	return httpRequest("POST", url, auth, connToken, httpRawDataBody(bodyRawData))
}

/*
Post Generic POST request
*/
func Post(url string, auth bool, connToken string, bodyJSON interface{}) (int, []byte, error) {
	return httpRequest("POST", url, auth, connToken, httpJSONBody(bodyJSON))
}

/*
PostStruct POST request that returns a struct of type 'map[string]interface{}'
*/
func PostStruct(url string, auth bool, connToken string, bodyJSON interface{}) (int, map[string]interface{}, error) {
	logs.GetLogger().Info(pathLOG + "[PostStruct] POST request [" + url + "] ...")

	status, data, err := Post(url, auth, connToken, bodyJSON)
	if err != nil {
		logs.GetLogger().Error(pathLOG+"[PostStruct] ERROR (1)", err)
		return status, nil, err
	}

	// create json
	var objmap map[string]interface{}
	if err := json.Unmarshal(data, &objmap); err != nil {
		logs.GetLogger().Error(pathLOG+"[PostStruct] ERROR (2)", err)
		return status, nil, err
	}

	return status, objmap, nil
}

///////////////////////////////////////////////////////////////////////////////
// DELETE

/*
Delete Generic DELETE request
*/
func Delete(url string, auth bool, connToken string, bodyJSON interface{}) (int, []byte, error) {
	return httpRequest("DELETE", url, auth, connToken, httpJSONBody(bodyJSON))
}

/*
DeleteStruct DELETE request that returns a struct of type 'map[string]interface{}'
*/
func DeleteStruct(url string, auth bool, connToken string) (int, map[string]interface{}, error) {
	logs.GetLogger().Info(pathLOG + "[DeleteStruct] DELETE request [" + url + "] ...")

	type Body struct {
		Content interface{}
	}

	status, data, err := Delete(url, auth, connToken, Body{})
	if err != nil {
		logs.GetLogger().Error(pathLOG+"[DeleteStruct] ERROR (1)", err)
		return status, nil, err
	}

	// create json
	var objmap map[string]interface{}
	if err := json.Unmarshal(data, &objmap); err != nil {
		logs.GetLogger().Error(pathLOG+"[DeleteStruct] WARNING (1)", err)
		return status, nil, err
	}

	return status, objmap, nil
}

///////////////////////////////////////////////////////////////////////////////
// PUT

/*
Put Generic PUT request
*/
func Put(url string, auth bool, connToken string, bodyJSON interface{}) (int, []byte, error) {
	return httpRequest("PUT", url, auth, connToken, httpJSONBody(bodyJSON))
}

/*
PutStruct PUT request that returns a struct of type 'map[string]interface{}'
*/
func PutStruct(url string, auth bool, connToken string, bodyJSON interface{}) (int, map[string]interface{}, error) {
	logs.GetLogger().Info(pathLOG + "[PutStruct] PUT request [" + url + "] ...")

	status, data, err := Put(url, auth, connToken, bodyJSON)
	if err != nil {
		logs.GetLogger().Error(pathLOG+"[PutStruct] ERROR (1)", err)
		return status, nil, err
	}

	// create json
	var objmap map[string]interface{}
	if err := json.Unmarshal(data, &objmap); err != nil {
		logs.GetLogger().Error(pathLOG+"[PutStruct] ERROR (2)", err)
		return status, nil, err
	}

	return status, objmap, nil
}

///////////////////////////////////////////////////////////////////////////////
// PATCH

/*
Patch Generic PATCH request
*/
func Patch(url string, auth bool, connToken string, bodyJSON interface{}) (int, []byte, error) {
	return httpRequest("PATCH", url, auth, connToken, httpJSONBody(bodyJSON))
}

/*
PatchStruct PATCH request that returns a struct of type 'map[string]interface{}'
*/
func PatchStruct(url string, auth bool, connToken string, bodyJSON interface{}) (int, map[string]interface{}, error) {
	logs.GetLogger().Info(pathLOG + "[PatchStruct] PATCH request [" + url + "] ...")

	status, data, err := Patch(url, auth, connToken, bodyJSON)
	if err != nil {
		logs.GetLogger().Error(pathLOG+"[PatchStruct] ERROR (1)", err)
		return status, nil, err
	}

	// create json
	var objmap map[string]interface{}
	if err := json.Unmarshal(data, &objmap); err != nil {
		logs.GetLogger().Error(pathLOG+"[PatchStruct] ERROR (2)", err)
		return status, nil, err
	}

	return status, objmap, nil
}
