package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/task-6/internal/db"
)

var errMock = errors.New("mock error")

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mockBehavior  func(mock sqlmock.Sqlmock)
		expectedNames []string
		expectError   bool
		errorContains string
	}{
		{
			name: "Success",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedNames: []string{"Alice", "Bob"},
			expectError:   false,
		},
		{
			name: "Query Error",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnError(errMock)
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "db query",
		},
		{
			name: "Rows Scan Error",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name", "age"}).
					AddRow("Alice", 25)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "rows scanning",
		},
		{
			name: "Rows Iteration Error",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					RowError(0, errMock)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "rows error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := db.New(mockDB)

			tc.mockBehavior(mock)

			names, err := service.GetNames()

			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorContains)
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

	tests := []struct {
		name          string
		mockBehavior  func(mock sqlmock.Sqlmock)
		expectedNames []string
		expectError   bool
		errorContains string
	}{
		{
			name: "Success",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedNames: []string{"Alice"},
			expectError:   false,
		},
		{
			name: "Query Error",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnError(errMock)
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "db query",
		},
		{
			name: "Rows Scan Error",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name", "extra"}).AddRow("Alice", "Extra")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "rows scanning",
		},
		{
			name: "Rows Iteration Error",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					RowError(0, errMock)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "rows error",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := db.New(mockDB)

			tc.mockBehavior(mock)

			names, err := service.GetUniqueNames()

			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
