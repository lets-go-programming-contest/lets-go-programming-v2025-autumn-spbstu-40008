package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/task-6/internal/db"
)

var mockDBErr = errors.New("simulated database error")

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		mockSetup     func(sqlmock.Sqlmock)
		expectedNames []string
		expectError   bool
		errSubstring  string
	}{
		{
			name: "returns list of names successfully",
			mockSetup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob")
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedNames: []string{"Alice", "Bob"},
			expectError:   false,
		},
		{
			name: "query fails",
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT name FROM users").
					WillReturnError(mockDBErr)
			},
			expectedNames: nil,
			expectError:   true,
			errSubstring:  "query execution failed",
		},
		{
			name: "scan fails due to column mismatch",
			mockSetup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name", "age"}).
					AddRow("Carol", 30)
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
			errSubstring:  "failed to read row data",
		},
		{
			name: "rows.Next succeeds but rows.Err returns error",
			mockSetup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Dave").
					RowError(0, mockDBErr)
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
			errSubstring:  "error occurred during iteration",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := db.New(mockDB)
			tc.mockSetup(mock)

			names, err := service.GetNames()

			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errSubstring)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		mockSetup     func(sqlmock.Sqlmock)
		expectedNames []string
		expectError   bool
		errSubstring  string
	}{
		{
			name: "returns distinct names",
			mockSetup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Eve")
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedNames: []string{"Eve"},
			expectError:   false,
		},
		{
			name: "distinct query fails",
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnError(mockDBErr)
			},
			expectedNames: nil,
			expectError:   true,
			errSubstring:  "distinct query failed",
		},
		{
			name: "distinct scan error",
			mockSetup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name", "extra"}).AddRow("Frank", "x")
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
			errSubstring:  "distinct row scan error",
		},
		{
			name: "distinct rows iteration succeeds but rows.Err fails",
			mockSetup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Grace").
					RowError(0, mockDBErr)
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
			errSubstring:  "distinct iteration error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := db.New(mockDB)
			tc.mockSetup(mock)

			names, err := service.GetUniqueNames()

			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errSubstring)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
