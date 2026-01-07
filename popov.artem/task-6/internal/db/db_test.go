package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/temaPop1e/lets-go-programming-v2025-autumn-spbstu-40008/popov.artem/task-6/internal/db"
)

var (
	errDatabase = errors.New("database error")
	errRow      = errors.New("row error")
)

type mockDB struct {
	queryFunc func(query string, args ...any) (*sql.Rows, error)
}

func (m *mockDB) Query(query string, args ...any) (*sql.Rows, error) {
	return m.queryFunc(query, args...)
}

func TestDataHandler_GetNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupMock      func(mock sqlmock.Sqlmock)
		expectedResult []string
		expectError    bool
		errorSubstring string
	}{
		{
			name: "success with multiple rows",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					AddRow("Charlie")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedResult: []string{"Alice", "Bob", "Charlie"},
		},
		{
			name: "error on query execution",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").WillReturnError(errDatabase)
			},
			expectError:    true,
			errorSubstring: "database query failed",
		},
		{
			name: "no records found",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectError:    true,
			errorSubstring: "no records found",
		},
		{
			name: "error on rows scan",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow(nil).
					AddRow("Charlie")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectError:    true,
			errorSubstring: "row processing error",
		},
		{
			name: "error on rows iteration",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					RowError(1, errRow)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectError:    true,
			errorSubstring: "row processing error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dbConn, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbConn.Close()

			handler := db.New(dbConn)

			tc.setupMock(mock)

			result, err := handler.GetNames()

			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorSubstring)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDataHandler_GetUniqueNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupMock      func(mock sqlmock.Sqlmock)
		expectedResult []string
		expectError    bool
		errorSubstring string
	}{
		{
			name: "success with duplicates",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					AddRow("Alice").
					AddRow("Charlie").
					AddRow("Bob")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedResult: []string{"Alice", "Bob", "Charlie"},
		},
		{
			name: "error on query execution",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errDatabase)
			},
			expectError:    true,
			errorSubstring: "database query failed",
		},
		{
			name: "no records found",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectError:    true,
			errorSubstring: "no records found",
		},
		{
			name: "error on rows scan",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow(nil).
					AddRow("Bob")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectError:    true,
			errorSubstring: "row processing error",
		},
		{
			name: "error on rows iteration",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					RowError(1, errRow)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectError:    true,
			errorSubstring: "row processing error",
		},
		{
			name: "success with single empty name",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedResult: []string{""},
		},
		{
			name: "success with multiple empty names",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("").
					AddRow("").
					AddRow("")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedResult: []string{""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dbConn, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbConn.Close()

			handler := db.New(dbConn)

			tc.setupMock(mock)

			result, err := handler.GetUniqueNames()

			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorSubstring)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDataHandler_New(t *testing.T) {
	t.Parallel()

	mockDB := &mockDB{
		queryFunc: func(query string, args ...any) (*sql.Rows, error) {
			return nil, errDatabase
		},
	}

	handler := db.New(mockDB)
	assert.NotNil(t, handler)
	assert.Equal(t, mockDB, handler.DB)
}