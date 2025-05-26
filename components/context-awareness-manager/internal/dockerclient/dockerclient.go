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
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

// DockerClient defines the interface for interacting with Docker.
type DockerClient interface {
	RunContainer(string, []string) (string, error)
}

// DockerConnector implements the DockerClient interface for Docker clients
type DockerConnector struct {
	client client.APIClient
}

// NewDockerConnector initializes a new DockerConnector with a Docker client.
func NewDockerConnector() (DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.WithHost("unix:///var/run/docker.sock"))
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %v", err)
	}
	logrus.Infof("Docker client created successfully")
	return &DockerConnector{client: cli}, nil
}

// getFullImageTag constructs the full image tag, ensuring the tag and registry are included.
func getFullImageTag(image string) string {
	// If the image doesn't have a tag, append ":latest"
	if !strings.Contains(image, ":") {
		image = fmt.Sprintf("%s:latest", image)
	}

	// Prepend the registry if available
	dockerServer := os.Getenv("DOCKER_SERVER")
	if dockerServer != "" {
		image = fmt.Sprintf("%s/%s", dockerServer, image)
	}

	return image
}

// ImageExistsLocally checks if an image exists locally.
func (docker *DockerConnector) ImageExistsLocally(ctx context.Context, imageName string) (bool, error) {
	images, err := docker.client.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("error listing images: %v", err)
	}

	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == imageName {
				return true, nil
			}
		}
	}
	return false, nil
}

// PullImage pulls the specified Docker image.
func (docker *DockerConnector) PullImage(ctx context.Context, imageName string) error {
	// Search image locally
	exists, err := docker.ImageExistsLocally(ctx, imageName)
	if err != nil {
		return fmt.Errorf("error verifying image locally: %v", err)
	}
	if exists {
		return nil
	}

	logrus.Infof("Image %s not found locally. Pulling...\n", imageName)
	out, err := docker.client.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("error pulling image: %v", err)
	}
	defer out.Close()

	type DockerPullStatus struct {
		Status   string `json:"status"`
		Progress string `json:"progress"`
		ID       string `json:"id"`
	}

	decoder := json.NewDecoder(out)
	for {
		var status DockerPullStatus
		if err := decoder.Decode(&status); err == io.EOF {
			break
		} else if err != nil {
			logrus.Warnf("Error decoding pull status: %v", err)
			break
		}
		if status.ID != "" {
			logrus.Infof("[%s] %s %s", status.ID, status.Status, status.Progress)
		} else {
			logrus.Infof("%s", status.Status)
		}
	}

	return nil
}

// CreateAndStartContainer creates and starts a Docker container from the specified image.
func (docker *DockerConnector) CreateAndStartContainer(ctx context.Context, imageName string, cmd []string) (string, error) {
	// Define the container configuration
	config := &container.Config{
		Image: imageName,
		Env: []string{
			"AGENT_ID=" + os.Getenv("AGENT_ID"),
		},
	}
	if len(cmd) > 0 {
		config.Cmd = cmd
	}

	// Create the container
	resp, err := docker.client.ContainerCreate(ctx, config, nil, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("error creating container: %v", err)
	}

	// Start the container
	if err := docker.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("error starting container: %v", err)
	}

	// Wait for the container to finish
	statusCh, errCh := docker.client.ContainerWait(ctx, resp.ID, "")
	select {
	case err := <-errCh:
		if err != nil {
			return "", fmt.Errorf("error waiting for container to finish: %v", err)
		}
	case <-statusCh:
	}

	return resp.ID, nil
}

// GetContainerLogs fetches logs from a running Docker container.
func (docker *DockerConnector) GetContainerLogs(ctx context.Context, containerID string) (string, error) {
	logOut, err := docker.client.ContainerLogs(ctx, containerID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", fmt.Errorf("error getting container logs: %v", err)
	}

	// Read logs into a buffer
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(logOut); err != nil {
		return "", fmt.Errorf("error reading container logs: %v", err)
	}

	logs := decodeDockerLogs(&buf)
	return logs, nil
}

// RunContainer deploys a container with the specified image and command
// and returns the container logs.
func (docker *DockerConnector) RunContainer(imageName string, cmd []string) (string, error) {
	// Create a context for Docker API calls
	ctx := context.Background()

	// Get the full image tag (with registry and tag)
	imageTag := getFullImageTag(imageName)
	// Pull the Docker image
	if err := docker.PullImage(ctx, imageTag); err != nil {
		return "", err
	}

	// Create and start the container
	containerID, err := docker.CreateAndStartContainer(ctx, imageTag, cmd)
	if err != nil {
		return "", err
	}

	// Get the container logs
	logs, err := docker.GetContainerLogs(ctx, containerID)
	if err != nil {
		return "", err
	}

	// Remove the container (forced)
	if err := docker.client.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true}); err != nil {
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
	return strings.TrimSpace(logs)
}
