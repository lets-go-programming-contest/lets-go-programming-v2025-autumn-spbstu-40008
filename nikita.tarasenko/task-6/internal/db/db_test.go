package db_test

import (
	"context"
	"errors"
	"testing"

	dataStorage "task-6/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	connectErr    = errors.New("connection failure")
	rowProcessErr = errors.New("row processing error")
	dbErr         = errors.New("database error")
)

func TestUserHandler_GetAllNames(t *testing.T) {
	t.Parallel()

	testScenarios := []struct {
		scenarioName string
		mockSetup    func(mock sqlmock.Sqlmock)
		expectError  bool
		expectedData []string
	}{
		{
			scenarioName: "successful name retrieval",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alexander").
					AddRow("Dmitry").
					AddRow("Sergey")
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rows)
			},
			expectError:  false,
			expectedData: []string{"Alexander", "Dmitry", "Sergey"},
		},
		{
			scenarioName: "query execution error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnError(connectErr)
			},
			expectError:  true,
			expectedData: nil,
		},
		{
			scenarioName: "data extraction error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alexander").
					AddRow(nil).
					AddRow("Sergey")
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rows)
			},
			expectError:  true,
			expectedData: nil,
		},
		{
			scenarioName: "row iteration error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alexander").
					RowError(0, rowProcessErr)
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rows)
			},
			expectError:  true,
			expectedData: nil,
		},
		{
			scenarioName: "empty result set",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rows)
			},
			expectError:  false,
			expectedData: nil,
		},
	}

	for _, ts := range testScenarios {
		t.Run(ts.scenarioName, func(t *testing.T) {
			t.Parallel()

			db, mock, createErr := sqlmock.New()
			require.NoError(t, createErr)
			defer db.Close()

			ts.mockSetup(mock)

			handler := dataStorage.InitializeUserHandler(db)

			result, err := handler.GetAllNames(context.Background())

			if ts.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				if ts.expectedData == nil {
					assert.Nil(t, result)
				} else {
					assert.Equal(t, ts.expectedData, result)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserHandler_GetDistinctNamesAsSet(t *testing.T) {
	t.Parallel()

	testScenarios := []struct {
		scenarioName string
		mockSetup    func(mock sqlmock.Sqlmock)
		expectError  bool
		expectedData []string
	}{
		{
			scenarioName: "successful distinct name retrieval",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alexander").
					AddRow("Alexander").
					AddRow("Dmitry").
					AddRow("Sergey").
					AddRow("Dmitry")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(rows)
			},
			expectError:  false,
			expectedData: []string{"Alexander", "Alexander", "Dmitry", "Sergey", "Dmitry"},
		},
		{
			scenarioName: "distinct query error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnError(dbErr)
			},
			expectError:  true,
			expectedData: nil,
		},
		{
			scenarioName: "distinct data extraction error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alexander").
					AddRow(nil).
					AddRow("Sergey")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(rows)
			},
			expectError:  true,
			expectedData: nil,
		},
		{
			scenarioName: "distinct row iteration error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alexander").
					RowError(0, rowProcessErr)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(rows)
			},
			expectError:  true,
			expectedData: nil,
		},
		{
			scenarioName: "empty distinct result set",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(rows)
			},
			expectError:  false,
			expectedData: nil,
		},
	}

	for _, ts := range testScenarios {
		t.Run(ts.scenarioName, func(t *testing.T) {
			t.Parallel()

			db, mock, createErr := sqlmock.New()
			require.NoError(t, createErr)
			defer db.Close()

			ts.mockSetup(mock)

			handler := dataStorage.InitializeUserHandler(db)

			result, err := handler.GetDistinctNamesAsSet(context.Background())

			if ts.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				if ts.expectedData == nil {
					assert.Nil(t, result)
				} else {
					assert.Equal(t, ts.expectedData, result)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
