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
package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	t.Run("Existing Environment Variable", func(t *testing.T) {
		key := "TEST_KEY"
		expected := "test_value"
		os.Setenv(key, expected)
		defer os.Unsetenv(key)

		value := getEnv(key, "default_value")
		assert.Equal(t, expected, value, "Should return environment variable value if set")
	})

	t.Run("Non-Existent Environment Variable", func(t *testing.T) {
		key := "NON_EXISTENT_KEY"
		defaultValue := "default_value"

		value := getEnv(key, defaultValue)
		assert.Equal(t, defaultValue, value, "Should return default value if environment variable is not set")
	})
}

func TestEnvInitialization(t *testing.T) {
	t.Run("Environment Variable Initialization", func(t *testing.T) {
		os.Setenv("ZENOH_URL", "http://test-url")
		os.Setenv("SERVER_PORT", "9090")
		defer os.Unsetenv("ZENOH_URL")
		defer os.Unsetenv("SERVER_PORT")

		zenohURL := getEnv("ZENOH_URL", "http://default-url")
		port := getEnv("SERVER_PORT", "8080")

		assert.Equal(t, "http://test-url", zenohURL, "ZENOH_URL should be set correctly")
		assert.Equal(t, "9090", port, "SERVER_PORT should be set correctly")
	})
}
