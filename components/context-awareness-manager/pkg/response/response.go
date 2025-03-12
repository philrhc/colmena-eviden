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
package response

import (
	"encoding/json"
	"log"
	"net/http"
)

// ERROR sends a JSON error response with the specified status code and error message.
// If err is nil, it sends a default bad request response.
func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(w, statusCode, map[string]string{"error": err.Error()})
		return
	}
	JSON(w, http.StatusBadRequest, map[string]string{"error": "Bad request"})
}

// JSON sends a JSON response with the specified status code and data.
// If encoding the data fails, it sends a 500 internal server error response.
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Only encode if there's data to encode
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Printf("Failed to encode JSON response: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

// Success is a helper function for sending successful responses
func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, data)
}

// InternalError is a helper function for sending internal server error responses
func InternalError(w http.ResponseWriter, err error) {
	ERROR(w, http.StatusInternalServerError, err)
}
