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

// Implementamos los métodos de client.APIClient que necesitamos
func (m *MockAPIClient) ImageList(ctx context.Context, opts image.ListOptions) ([]image.ListOptions, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]image.ListOptions), args.Error(1)
}

func (m *MockAPIClient) ImagePull(ctx context.Context, ref string, options image.PullOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, ref, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func TestGetFullImageTag(t *testing.T) {
	// Sin configurar el servidor Docker
	result := getFullImageTag("test-image")
	assert.Equal(t, "test-image:latest", result)

	// Configurando la variable DOCKER_SERVER
	os.Setenv("DOCKER_SERVER", "mydockerregistry.com")
	defer os.Unsetenv("DOCKER_SERVER")

	result = getFullImageTag("test-image")
	assert.Equal(t, "mydockerregistry.com/test-image:latest", result)
}

func TestDecodeDockerLogs(t *testing.T) {
	// Simulamos un log de Docker que consiste en un encabezado de 8 bytes y un payload
	// El encabezado tiene el siguiente formato: [00 00 00 00 00 00 00 10] (payload de 16 bytes)
	header := []byte{0, 0, 0, 0, 0, 0, 0, 23} // 23 bytes de payload
	payload := []byte("This is the log message")
	logData := append(header, payload...)

	// Creamos un reader con los datos simulados
	reader := bytes.NewReader(logData)

	// Llamamos a la función a probar
	result := decodeDockerLogs(reader)

	// Comprobamos que el resultado es el esperado
	expected := "This is the log message"
	assert.Equal(t, expected, result)
}
