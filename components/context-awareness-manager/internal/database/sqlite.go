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
package database

import (
	"context-awareness-manager/pkg/models"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // Driver SQLite
	"github.com/sirupsen/logrus"
)

// Database defines the interface for database operations
type Database interface {
	Close()
	InsertContext(models.DockerContextDefinition) error
	InsertContexts([]models.DockerContextDefinition) error
	DeleteContext(string) error
	GetContextByID(string) (*models.DockerContextDefinition, error)
	GetAllContexts() ([]models.DockerContextDefinition, error)
	InsertClassification(string, string) error
	GetLastContextClassification(string) (string, error)
}

// SQLConnector implements the Database interface
type SQLConnector struct {
	DB *sql.DB
}

// NewSQLConnector initializes the database and returns an SQLConnector instance.
func NewSQLConnector(dbPath string) (*SQLConnector, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	connector := &SQLConnector{DB: db}

	// Create Context Table if it does not exist
	err = connector.CreateTables()
	if err != nil {
		return nil, fmt.Errorf("error creating tables: %v", err)
	}

	logrus.Info("Database initialized successfully")
	return connector, nil
}

// CloseDB closes the database connection
func (s *SQLConnector) Close() {
	if s.DB != nil {
		s.DB.Close()
	}
}

// createTables ensures required tables exist in the database.
func (s *SQLConnector) CreateTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS dockerContextDefinitions (
			id TEXT PRIMARY KEY,
			imageId TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS dockerContextClassifications (
			imageId TEXT PRIMARY KEY,
			classification TEXT NOT NULL
		)`,
	}

	for _, query := range queries {
		if _, err := s.DB.Exec(query); err != nil {
			return fmt.Errorf("error creating table: %w", err)
		}
	}

	logrus.Info("Database tables verified/created successfully")
	return nil
}

// InsertContext inserts a single context into the database
func (s *SQLConnector) InsertContext(context models.DockerContextDefinition) error {
	_, err := s.DB.Exec(`INSERT INTO dockerContextDefinitions (id, imageId) VALUES (?, ?)`,
		context.ID, context.ImageID)
	if err != nil {
		return fmt.Errorf("failed to insert/update context ID %s: %v", context.ID, err)
	}
	return nil
}

// InsertContexts inserts multiple context in a single transaction into the database
func (s *SQLConnector) InsertContexts(contexts []models.DockerContextDefinition) error {
	// Begin a transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback() // Ensure the transaction is rolled back if something goes wrong

	// Prepare the insert statement for batch insert
	stmt, err := tx.Prepare("INSERT OR REPLACE INTO DockerContextDefinitions (id, imageId) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %v", err)
	}
	defer stmt.Close()

	// Insert each context into the database
	for _, context := range contexts {
		_, err := stmt.Exec(context.ID, context.ImageID)
		if err != nil {
			return fmt.Errorf("failed to execute insert for context %s: %v", context.ID, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// DeleteContext deletes a context from the database
func (s *SQLConnector) DeleteContext(contextID string) error {
	_, err := s.DB.Exec(`DELETE FROM dockerContextDefinitions WHERE id = ?`, contextID)
	if err != nil {
		return fmt.Errorf("failed to delete context ID %s: %v", contextID, err)
	}

	logrus.Debugf("Context ID %s deleted successfully into database", contextID)
	return nil
}

// GetContextByID retrieves a context by its ID from the database
func (s *SQLConnector) GetContextByID(contextID string) (*models.DockerContextDefinition, error) {
	var context models.DockerContextDefinition
	err := s.DB.QueryRow(`SELECT id, imageId FROM dockerContextDefinitions WHERE id = ?`, contextID).
		Scan(&context.ID, &context.ImageID)

	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Warnf("No context found for ID %s", contextID)
			return nil, nil // Context not found
		}
		return nil, fmt.Errorf("failed to get context ID %s: %v", contextID, err)
	}

	logrus.Debugf("Context ID %s retrieved successfully", contextID)
	return &context, nil
}

// GetAllContexts retrieves all stored contexts from the database
func (s *SQLConnector) GetAllContexts() ([]models.DockerContextDefinition, error) {
	rows, err := s.DB.Query(`SELECT id, imageId FROM dockerContextDefinitions`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all contexts: %v", err)
	}
	defer rows.Close()

	var contexts []models.DockerContextDefinition
	for rows.Next() {
		var context models.DockerContextDefinition
		if err := rows.Scan(&context.ID, &context.ImageID); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		contexts = append(contexts, context)
	}

	return contexts, nil
}

// InsertContext inserts a single context into the database
func (s *SQLConnector) InsertClassification(imageID string, classification string) error {
	_, err := s.DB.Exec(`INSERT INTO dockerContextClassifications (imageId, classification) VALUES (?, ?)`,
		imageID, classification)
	if err != nil {
		return fmt.Errorf("failed to insert/update classification for %s: %v", imageID, err)
	}
	return nil
}

// GetLastContextClassification retrieves all stored contexts from the database
func (s *SQLConnector) GetLastContextClassification(imageID string) (string, error) {
	var classification string
	err := s.DB.QueryRow(`SELECT classification FROM dockerContextClassifications WHERE imageId = ?`, imageID).
		Scan(&classification)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Context not found
		}
		return "", fmt.Errorf("failed to get classification for %s: %v", imageID, err)
	}
	return classification, nil
}
