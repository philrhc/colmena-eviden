package database

import (
	"context-awareness-manager/pkg/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestInsertContext(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create SQLConnector using the mocked database
	connector := SQLConnector{DB: db}

	// Define the context you want to insert
	context := models.DockerContextDefinition{
		ID:      "context1",
		ImageID: "image123",
	}

	// Set up the expected query and response
	mock.ExpectExec("INSERT INTO dockerContextDefinitions").
		WithArgs(context.ID, context.ImageID).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Mock the result of the query

	// Call the method under test
	err = connector.InsertContext(context)
	assert.NoError(t, err)

	// Ensure that all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertContexts(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create SQLConnector using the mocked database
	connector := SQLConnector{DB: db}

	// Define a list of contexts to insert
	contexts := []models.DockerContextDefinition{
		{ID: "context1", ImageID: "image123"},
		{ID: "context2", ImageID: "image456"},
	}

	// Start a transaction for batch insert
	mock.ExpectBegin()

	// Set up expected statements for each insert
	stmt := mock.ExpectPrepare("INSERT OR REPLACE INTO DockerContextDefinitions")
	stmt.ExpectExec().
		WithArgs(contexts[0].ID, contexts[0].ImageID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	stmt.ExpectExec().
		WithArgs(contexts[1].ID, contexts[1].ImageID).
		WillReturnResult(sqlmock.NewResult(2, 1))

	// Expect commit after the batch insert
	mock.ExpectCommit()

	// Call the method under test
	err = connector.InsertContexts(contexts)
	assert.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetContextByID(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create SQLConnector using the mocked database
	connector := SQLConnector{DB: db}

	// Define the context ID to retrieve
	contextID := "context1"

	// Define the mock result
	rows := sqlmock.NewRows([]string{"id", "imageId"}).
		AddRow("context1", "image123")

	// Set up the expected query and rows returned
	mock.ExpectQuery("SELECT id, imageId FROM dockerContextDefinitions").
		WithArgs(contextID).
		WillReturnRows(rows)

	// Call the method under test
	context, err := connector.GetContextByID(contextID)
	assert.NoError(t, err)
	assert.NotNil(t, context)
	assert.Equal(t, "context1", context.ID)
	assert.Equal(t, "image123", context.ImageID)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteContext(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create SQLConnector using the mocked database
	connector := SQLConnector{DB: db}

	// Define the context ID to delete
	contextID := "context1"

	// Set up the expected query and result
	mock.ExpectExec("DELETE FROM dockerContextDefinitions").
		WithArgs(contextID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the method under test
	err = connector.DeleteContext(contextID)
	assert.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllContexts(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create SQLConnector using the mocked database
	connector := SQLConnector{DB: db}

	// Define the mock result
	rows := sqlmock.NewRows([]string{"id", "imageId"}).
		AddRow("context1", "image123").
		AddRow("context2", "image456")

	// Set up the expected query and rows returned
	mock.ExpectQuery("SELECT id, imageId FROM dockerContextDefinitions").
		WillReturnRows(rows)

	// Call the method under test
	contexts, err := connector.GetAllContexts()
	assert.NoError(t, err)
	assert.Len(t, contexts, 2)
	assert.Equal(t, "context1", contexts[0].ID)
	assert.Equal(t, "context2", contexts[1].ID)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertClassification(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create SQLConnector using the mocked database
	connector := SQLConnector{DB: db}

	// Define the image ID and classification to insert
	imageID := "image123"
	classification := "critical"

	// Set up the expected query and response
	mock.ExpectExec("INSERT OR REPLACE INTO dockerContextClassifications").
		WithArgs(imageID, classification).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the method under test
	err = connector.InsertClassification(imageID, classification)
	assert.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLastContextClassification(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create SQLConnector using the mocked database
	connector := SQLConnector{DB: db}

	// Define the image ID to retrieve the classification
	imageID := "image123"

	// Define the mock result
	rows := sqlmock.NewRows([]string{"classification"}).
		AddRow("critical")

	// Set up the expected query and rows returned
	mock.ExpectQuery("SELECT classification FROM dockerContextClassifications").
		WithArgs(imageID).
		WillReturnRows(rows)

	// Call the method under test
	classification, err := connector.GetLastContextClassification(imageID)
	assert.NoError(t, err)
	assert.Equal(t, "critical", classification)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
