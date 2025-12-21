package db

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errDBFailure = errors.New("database connection failed")
	errRowFailure = errors.New("corrupted row data")
)

func TestDataHandler_RetrieveNames(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		mockSetup func(sqlmock.Sqlmock)
		expected []string
		expectError bool
		errorContains string
	}{
		{
			name: "valid multiple records",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("User1").
					AddRow("User2").
					AddRow("User3")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expected: []string{"User1", "User2", "User3"},
		},
		{
			name: "database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").WillReturnError(errDBFailure)
			},
			expectError: true,
			errorContains: "database query failed",
		},
		{
			name: "empty result",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectError: true,
			errorContains: "no records found",
		},
		{
			name: "row processing error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(123)
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
			handler := CreateHandler(mockDB)
			tc.mockSetup(mock)
			result, err := handler.RetrieveNames()
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

func TestDataHandler_RetrieveUniqueNames(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		mockSetup func(sqlmock.Sqlmock)
		expected []string
		expectError bool
		errorContains string
	}{
		{
			name: "valid distinct records",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("UserA").
					AddRow("UserB").
					AddRow("UserA")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expected: []string{"UserA", "UserB"},
		},
		{
			name: "query execution error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(sql.ErrNoRows)
			},
			expectError: true,
			errorContains: "database query failed",
		},
		{
			name: "no distinct records",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectError: true,
			errorContains: "no distinct records",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()
			handler := CreateHandler(mockDB)
			tc.mockSetup(mock)
			result, err := handler.RetrieveUniqueNames()
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.ElementsMatch(t, tc.expected, result)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
