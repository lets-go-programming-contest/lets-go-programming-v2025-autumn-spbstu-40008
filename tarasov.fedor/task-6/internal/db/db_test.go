package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBService_GetNames(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	service := New(db)

	tests := []struct {
		name          string
		mockBehavior  func()
		expectedNames []string
		expectError   bool
		errorContains string
	}{
		{
			name: "Success",
			mockBehavior: func() {
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
			mockBehavior: func() {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnError(errors.New("connection refused"))
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "db query",
		},
		{
			name: "Rows Scan Error",
			mockBehavior: func() {
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
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					RowError(0, errors.New("row corrupted"))
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "rows error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			names, err := service.GetNames()

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedNames, names)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	service := New(db)

	tests := []struct {
		name          string
		mockBehavior  func()
		expectedNames []string
		expectError   bool
		errorContains string
	}{
		{
			name: "Success",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedNames: []string{"Alice"},
			expectError:   false,
		},
		{
			name: "Query Error",
			mockBehavior: func() {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnError(errors.New("fail"))
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "db query",
		},
		{
			name: "Rows Scan Error",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"name", "extra"}).AddRow("Alice", "Extra")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "rows scanning",
		},
		{
			name: "Rows Iteration Error",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					RowError(0, errors.New("iteration error"))
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "rows error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()
			names, err := service.GetUniqueNames()
			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedNames, names)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
