package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Czeeen/lets-go-programming-v2025-autumn-spbstu-40008/nikita.prokopenko/task-6/internal/db"
)

var (
	errDBFailure = errors.New("database connection failed")
	errRowFailure = errors.New("corrupted row data")
)

func TestDataHandler_GetNames(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		setupMock func(sqlmock.Sqlmock)
		expected []string
		expectError bool
		errorContains string
	}{
		{
			name: "success - return names",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					AddRow("Charlie")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expected: []string{"Alice", "Bob", "Charlie"},
		},
		{
			name: "error - query error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").WillReturnError(errDBFailure)
			},
			expectError: true,
			errorContains: "database query failed",
		},
		{
			name: "success - empty result",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectError: true,
			errorContains: "no records found",
		},
		{
			name: "error - scan error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectError: true,
			errorContains: "row processing error",
		},
		{
			name: "error - rows.Err",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					RowError(1, errRowFailure)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectError: true,
			errorContains: "row processing error",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()
			handler := db.New(mockDB)
			tc.setupMock(mock)
			result, err := handler.GetNames()
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDataHandler_GetUniqueNames(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		setupMock func(sqlmock.Sqlmock)
		expected []string
		expectError bool
		errorContains string
	}{
		{
			name: "success - return values",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expected: []string{"Alice", "Bob"},
		},
		{
			name: "error - query error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
			errorContains: "database query failed",
		},
		{
			name: "error - scan error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectError: true,
			errorContains: "row processing error",
		},
		{
			name: "error - rows.Err",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					RowError(1, errRowFailure)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectError: true,
			errorContains: "row processing error",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()
			handler := db.New(mockDB)
			tc.setupMock(mock)
			result, err := handler.GetUniqueNames()
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
