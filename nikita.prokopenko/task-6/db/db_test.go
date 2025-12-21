package usersdb_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Czeeen/lets-go-programming-v2025-autumn-spbstu-40008/nikita.prokopenko/task-6/internal/usersdb"
)

var (
	errDBConnection = errors.New("database connection error")
	errRowCorrupted = errors.New("corrupted row data")
)

func TestUserDataService_FetchUsernames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupMock     func(sqlmock.Sqlmock)
		expectedNames []string
		expectError   bool
		errorPart     string
	}{
		{
			name: "successful fetch with multiple users",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alex").
					AddRow("Maria").
					AddRow("Ivan")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedNames: []string{"Alex", "Maria", "Ivan"},
		},
		{
			name: "database query failure",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").WillReturnError(errDBConnection)
			},
			expectError: true,
			errorPart:   "database query execution failed",
		},
		{
			name: "empty result set",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectError: true,
			errorPart:   "no rows available in result set",
		},
		{
			name: "row scanning error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(123)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectError: true,
			errorPart:   "error scanning row data",
		},
		{
			name: "result processing error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alex").
					AddRow("Maria").
					RowError(1, errRowCorrupted)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectError: true,
			errorPart:   "error processing query results",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := usersdb.NewUserService(mockDB)
			tt.setupMock(mock)

			names, err := service.FetchUsernames()

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorPart)
				assert.Nil(t, names)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedNames, names)
			}
			
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserDataService_FetchUniqueUsernames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupMock     func(sqlmock.Sqlmock)
		expectedNames []string
		expectError   bool
		errorPart     string
	}{
		{
			name: "successful fetch with unique names",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alex").
					AddRow("Maria").
					AddRow("Alex")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedNames: []string{"Alex", "Maria"},
		},
		{
			name: "query execution error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
			errorPart:   "database query execution failed",
		},
		{
			name: "no unique records",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectError: true,
			errorPart:   "no unique records found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := usersdb.NewUserService(mockDB)
			tt.setupMock(mock)

			names, err := service.FetchUniqueUsernames()

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorPart)
				assert.Nil(t, names)
			} else {
				require.NoError(t, err)
				assert.ElementsMatch(t, tt.expectedNames, names)
			}
			
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestNewUserService(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := usersdb.NewUserService(mockDB)

	require.NotNil(t, service)
	assert.Equal(t, mockDB, service.Executor())
}

func (s *usersdb.UserDataService) Executor() usersdb.QueryExecutor {
	return s.executor
}
