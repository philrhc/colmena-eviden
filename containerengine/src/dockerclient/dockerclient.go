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
package docker

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// RunContainer deploys a container with the specified image and command
// and returns the container logs
func RunContainer(image string, cmd []string) (string, error) {
	host := os.Getenv("DOCKER_HOST")
	if host == "" {
		host = "unix:///var/run/docker.sock"
	}
	// Create a Docker client
	cli, err := client.NewClientWithOpts(client.WithHost(host), client.WithAPIVersionNegotiation())
	if err != nil {
		return "", fmt.Errorf("error creating Docker client: %v", err)
	}

	// Define the container configuration
	config := &container.Config{
		Image: image,
	}
	if len(cmd) > 0 {
		config.Cmd = cmd
	}

	// Create the container
	resp, err := cli.ContainerCreate(context.Background(), config, nil, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("error creating container: %v", err)
	}

	// Start the container
	if err := cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("error starting container: %v", err)
	}

	fmt.Printf("Container started with ID: %s\n", resp.ID)

	// Wait for the container to finish
	statusCh, errCh := cli.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", fmt.Errorf("error waiting for container to finish: %v", err)
		}
	case <-statusCh:
	}

	// Get the container logs
	out, err := cli.ContainerLogs(context.Background(), resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", fmt.Errorf("error getting container logs: %v", err)
	}

	// Read logs into a buffer
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(out); err != nil {
		return "", fmt.Errorf("error reading container logs: %v", err)
	}

	logs := decodeDockerLogs(&buf)

	// Remove the container
	if err := cli.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{}); err != nil {
		return "", fmt.Errorf("error removing container: %v", err)
	}

	return logs, nil
}

// decodeDockerLogs decodes the Docker log format and returns plain text logs
func decodeDockerLogs(reader io.Reader) string {
	var logs string
	buf := make([]byte, 8) // Docker log header is 8 bytes
	for {
		_, err := io.ReadFull(reader, buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		// Docker log stream has an 8-byte header
		payloadLen := binary.BigEndian.Uint32(buf[4:8])
		payload := make([]byte, payloadLen)
		_, err = io.ReadFull(reader, payload)
		if err != nil {
			break
		}
		logs += string(payload)
	}
	return logs
}
