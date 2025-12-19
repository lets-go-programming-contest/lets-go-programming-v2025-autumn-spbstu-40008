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

var (
	errConnectionFailed = errors.New("connection failed")
	errRowError         = errors.New("row error")
)

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		expected    []string
		expectedErr string
	}{
		{
			name: "success - return names",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					AddRow("Charlie")

				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rows)
			},
			expected: []string{"Alice", "Bob", "Charlie"},
		},
		{
			name: "error - query error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnError(errConnectionFailed)
			},
			expectedErr: "db query:",
		},
		{
			name: "success - empty result",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})

				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rows)
			},
			expected: []string{},
		},
		{
			name: "error - scan error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow(123)

				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rows)
			},
			expectedErr: "rows scanning:",
		},
		{
			name: "error - rows.Err",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					RowError(1, errRowError)

				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rows)
			},
			expectedErr: "rows error:",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := db.New(mockDB)

			tc.setupMock(mock)

			result, err := service.GetNames()

			if tc.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
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
		expectedErr string
	}{
		{
			name: "success - return values",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob")

				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(rows)
			},
			expected: []string{"Alice", "Bob"},
		},
		{
			name: "error - query error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnError(sql.ErrConnDone)
			},
			expectedErr: "db query:",
		},
		{
			name: "error - scan error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow(123)

				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(rows)
			},
			expectedErr: "rows scanning:",
		},
		{
			name: "error - rows.Err",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					RowError(1, errRowError)

				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(rows)
			},
			expectedErr: "rows error:",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := db.New(mockDB)

			tc.setupMock(mock)

			result, err := service.GetUniqueNames()

			if tc.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetNames_SingleRow(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice")

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	service := db.New(mockDB)

	result, err := service.GetNames()

	require.NoError(t, err)
	assert.Equal(t, []string{"Alice"}, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames_RowsCloseError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	closeErr := errors.New("close error")

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		CloseError(closeErr)

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	service := db.New(mockDB)

	result, err := service.GetNames()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "rows close")
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
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
