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
package monitor

import (
	"context-awareness-manager/internal/database"
	"context-awareness-manager/internal/dockerclient"
	"context-awareness-manager/internal/publisher"
	"context-awareness-manager/pkg/models"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// ContextMonitor defines the interface for handling Docker contexts
type ContextMonitor interface {
	RegisterContexts([]models.DockerContextDefinition) error
	StartMonitoring()
}

// contextMonitorImpl implements the ContextMonitor interface.
type contextMonitorImpl struct {
	dbConnection *database.SQLConnector
	dockerClient dockerclient.DockerClient
	publisher    publisher.Publisher
}

// NewContextMonitor creates a new instance of contextMonitorImpl.
func NewContextMonitor(dbConnection *database.SQLConnector, dockerClient dockerclient.DockerClient, publisher publisher.Publisher) ContextMonitor {
	return &contextMonitorImpl{
		dbConnection: dbConnection,
		dockerClient: dockerClient,
		publisher:    publisher,
	}
}

// RegisterContexts inserts new Docker contexts into the database
func (c *contextMonitorImpl) RegisterContexts(newContexts []models.DockerContextDefinition) error {
	// Insert new contexts into the database
	err := c.dbConnection.InsertContexts(newContexts)
	if err != nil {
		logrus.Errorf("Error inserting new contexts into database: %v", err)
	}
	logrus.Infof("Successfully registered %d new contexts", len(newContexts))
	return nil
}

// StartContextMonitoring listens for new contexts and processes them dynamically.
func (c *contextMonitorImpl) StartMonitoring() {

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {

		// Get stored contexts from the database
		contexts, err := c.dbConnection.GetAllContexts()
		if err != nil {
			logrus.Errorf("Error retrieving contexts from database: %v", err)
		}

		if len(contexts) == 0 {
			logrus.Debug("No contexts to monitor, waiting for new ones...")
			continue
		}

		// Iterate over each context and classify + publish
		for _, ctx := range contexts {
			// Classify the context
			classification, error := c.ClassifyContext(ctx)
			if error != nil {
				logrus.Errorf("Error processing context %s: %v", ctx.ID, error)
				continue
			}

			// Construct the key expression for the context
			keyExpression := fmt.Sprintf("colmena/contexts/%s/%s", os.Getenv("AGENT_ID"), ctx.ID)

			if classification != "" {
				err := c.publisher.PublishContextClassification(keyExpression, classification)
				if err != nil {
					logrus.Errorf("Error publishing context: %v", err)
					continue
				}
				logrus.Info("Context published successfully!")
			}
		}
	}
}

// ProcessContextClassifications checks and updates the classification of each context
func (c *contextMonitorImpl) ClassifyContext(ctx models.DockerContextDefinition) (string, error) {
	// Run the container and obtain the classification
	classification, err := c.dockerClient.RunContainer(ctx.ImageID, []string{})
	if err != nil {
		logrus.Errorf("Error deploying container for context %s: %v\n", ctx.ID, err)
		return "", err
	}

	// Get the last classification for the context
	lastClassification, err := c.dbConnection.GetLastContextClassification(ctx.ImageID)
	if err != nil {
		logrus.Errorf("Error retrieving last classification for %s: %v\n", ctx.ID, err)
		return "", err
	}

	// If no previous classification exists, insert the new one directly
	if lastClassification == "" {
		logrus.Infof("No previous classification found for %s. Inserting initial classification: %s", ctx.ID, classification)
		err = c.dbConnection.InsertClassification(ctx.ImageID, classification)
		if err != nil {
			logrus.Errorf("Error saving classification for %s: %v", ctx.ID, err)
			return "", err
		}
		return classification, nil
	}

	// If classification hasn't changed, return empty string to indicate no change
	if classification == lastClassification {
		return "", nil
	}

	// If classification has changed, insert the new classification into the database
	logrus.Infof("Context classification for %s: %s\n", ctx.ID, classification)
	err = c.dbConnection.InsertClassification(ctx.ImageID, classification)
	if err != nil {
		logrus.Errorf("Error saving classification for %s: %v", ctx.ID, err)
		return "", err
	}

	// Return the new classification
	return classification, nil
}
