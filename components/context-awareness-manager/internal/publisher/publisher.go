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
package publisher

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// Publisher defines the interface for publishing and deleting context classifications.
type Publisher interface {
	PublishContextClassification(string, string) error
	DeleteContext(string) error
}

type zenohPublisher struct {
	client   *http.Client
	zenohURL string
}

// NewPublisher creates a new instance of the Zenoh-based publisher.
func NewPublisher(zenohURL string) Publisher {
	return &zenohPublisher{
		client:   &http.Client{},
		zenohURL: zenohURL,
	}
}

// sendHTTPRequest is a helper function to send HTTP requests and return the response.
func (z *zenohPublisher) sendHTTPRequest(method, url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating %s request: %v", method, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := z.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending %s request: %v", method, err)
	}

	return resp, nil
}

// PublishContextClassification sends an HTTP PUT request to the specified URL with a JSON body
// containing the provided value. The function takes two string parameters: `url` and `value`.
// It returns an error if the PUT request fails or if the response status code is not 200 OK.
func (z *zenohPublisher) PublishContextClassification(key string, value string) error {
	// Construct the URL using the provided key
	url := fmt.Sprintf("%s/%s", z.zenohURL, key)

	logrus.Infof("Sending PUT request to %s with body %s", url, value)
	resp, err := z.sendHTTPRequest(http.MethodPut, url, []byte(value))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// DeleteContext sends an HTTP DELETE request to the specified URL to remove a context.
// The function takes a `key` parameter representing the context to be deleted and returns an error if the DELETE request fails
// or if the response status code is not 200 OK.
func (z *zenohPublisher) DeleteContext(key string) error {
	// Construct the URL using the provided key
	url := fmt.Sprintf("%s/%s", z.zenohURL, key)

	logrus.Infof("Sending DELETE request to %s", url)
	resp, err := z.sendHTTPRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	return nil
}
