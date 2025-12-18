package db

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBService_GetNames(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	dbService := New(mockDB)

	testCases := []struct {
		name        string
		setupMock   func()
		expected    []string
		expectedErr error
	}{
		{
			name: "success - return names",
			setupMock: func() {
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
			setupMock: func() {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnError(errors.New("connection failed"))
			},
			expectedErr: errors.New("db query: connection failed"),
		},
		{
			name: "success - empty result",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expected: []string{},
		},
		{
			name: "error - rows reading error",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow(nil)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedErr: errors.New("rows scanning:"),
		},
		{
			name: "error - rows.Err()",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					RowError(1, errors.New("row error"))
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedErr: errors.New("rows error: row error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

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
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	dbService := New(mockDB)

	testCases := []struct {
		name        string
		setupMock   func()
		expected    []string
		expectedErr error
	}{
		{
			name: "success - return unique names",
			setupMock: func() {
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
			setupMock: func() {
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
			setupMock: func() {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnError(sql.ErrConnDone)
			},
			expectedErr: errors.New("db query:"),
		},
		{
			name: "success - empty result",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expected: []string{},
		},
		{
			name: "error - scan error",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow(123)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedErr: errors.New("rows scanning:"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

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
	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := New(mockDB)

	assert.NotNil(t, service)
	assert.Equal(t, mockDB, service.DB)
}
