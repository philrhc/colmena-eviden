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
package dockerclient

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/docker/docker/api/types/image"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAPIClient is a mock type for the Docker API client
type MockAPIClient struct {
	mock.Mock
}

// Implement the methods of client.APIClient that we need
func (m *MockAPIClient) ImageList(ctx context.Context, opts image.ListOptions) ([]image.ListOptions, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]image.ListOptions), args.Error(1)
}

func (m *MockAPIClient) ImagePull(ctx context.Context, ref string, options image.PullOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, ref, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func TestGetFullImageTag(t *testing.T) {
	// Without configuring the Docker server
	result := getFullImageTag("test-image")
	assert.Equal(t, "test-image:latest", result)

	// Configuring the DOCKER_SERVER variable
	os.Setenv("DOCKER_SERVER", "mydockerregistry.com")
	defer os.Unsetenv("DOCKER_SERVER")

	result = getFullImageTag("test-image")
	assert.Equal(t, "mydockerregistry.com/test-image:latest", result)
}

func TestDecodeDockerLogs(t *testing.T) {
	// Simulate a Docker log consisting of an 8-byte header and a payload
	// The header has the following format: [00 00 00 00 00 00 00 10] (16-byte payload)
	header := []byte{0, 0, 0, 0, 0, 0, 0, 23} // 23-byte payload
	payload := []byte("This is the log message")
	logData := append(header, payload...)

	// Create a reader with the simulated data
	reader := bytes.NewReader(logData)

	// Call the function to test
	result := decodeDockerLogs(reader)

	// Verify that the result is as expected
	expected := "This is the log message"
	assert.Equal(t, expected, result)
}
