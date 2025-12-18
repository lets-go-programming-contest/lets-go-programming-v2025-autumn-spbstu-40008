package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	db "github.com/julia.pshenitsyna/task-6/internal/db"
)

var(
	errConnectionFailed = errors.New("connection failed")
    errRowError         = errors.New("row error")
    errScanning         = errors.New("rows scanning")
    errRows             = errors.New("rows error")
)

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		expected    []string
		expectedErr error
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
			name: "error - error request",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnError(errConnectionFailed)
			},
			expectedErr: errors.New("db query: " + errConnectionFailed.Error()),
		},
		{
			name: "success - empty result",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expected: []string{},
		},
		{
			name: "error - rows reading error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow(nil)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedErr: errors.New("rows scanning: " + errScanning.Error()),
		},
		{
			name: "error - rows.Err()",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					RowError(1, errRowError)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedErr: errors.New("rows error: " + errRowError.Error()),
		},
		{
			name: "error - sql.ErrConnDone",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnError(sql.ErrConnDone)
			},
			expectedErr: errors.New("db query: " + sql.ErrConnDone.Error()),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			dbService := db.New(mockDB)
			tc.setupMock(mock)

			result, err := dbService.GetNames()

			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr.Error())
				
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		expected    []string
		expectedErr error
	}{
		{
			name: "success - return unique names",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					AddRow("Alice").
					AddRow("Charlie")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expected: []string{"Alice", "Bob", "Alice", "Charlie"},
		},
		{
			name: "success - only unique names",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					AddRow("Charlie")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expected: []string{"Alice", "Bob", "Charlie"},
		},
		{
			name: "error - request error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnError(sql.ErrConnDone)
			},
			expectedErr: errors.New("db query: " + sql.ErrConnDone.Error()),
		},
		{
			name: "success - empty result",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expected: []string{},
		},
		{
			name: "error - scan error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow(123)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedErr: errors.New("rows scanning: " + errScanning.Error()),
		},
		{
			name: "success - single name",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Single")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expected: []string{"Single"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			dbService := db.New(mockDB)
			tc.setupMock(mock)

			result, err := dbService.GetUniqueNames()

			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr.Error())
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()
	
	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	assert.NotNil(t, service)
	assert.Equal(t, mockDB, service.DB)
}
