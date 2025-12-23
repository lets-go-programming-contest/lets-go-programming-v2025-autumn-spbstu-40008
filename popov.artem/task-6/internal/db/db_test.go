package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/task-6/internal/db"
)

var simErr = errors.New("simulated database failure")

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		title         string
		setupMocks    func(sqlmock.Sqlmock)
		expectedNames []string
		shouldErr     bool
		errFragment   string
	}{
		{
			title: "successful retrieval",
			setupMocks: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT name FROM users").
					WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("UserA").AddRow("UserB"))
			},
			expectedNames: []string{"UserA", "UserB"},
			shouldErr:     false,
		},
		{
			title: "query error",
			setupMocks: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT name FROM users").
					WillReturnError(simErr)
			},
			expectedNames: nil,
			shouldErr:     true,
			errFragment:   "query execution failed",
		},
		{
			title: "scan mismatch",
			setupMocks: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT name FROM users").
					WillReturnRows(sqlmock.NewRows([]string{"name", "age"}).AddRow("Test", 100))
			},
			expectedNames: nil,
			shouldErr:     true,
			errFragment:   "failed to read row data",
		},
		{
			title: "row iteration error",
			setupMocks: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					RowError(0, simErr)
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			shouldErr:     true,
			errFragment:   "error occurred during iteration",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := db.New(mockDB)
			tc.setupMocks(mock)

			names, err := service.GetNames()

			if tc.shouldErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errFragment)
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
		title         string
		setupMocks    func(sqlmock.Sqlmock)
		expectedNames []string
		shouldErr     bool
		errFragment   string
	}{
		{
			title: "distinct success",
			setupMocks: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("OnlyOne"))
			},
			expectedNames: []string{"OnlyOne"},
			shouldErr:     false,
		},
		{
			title: "distinct query error",
			setupMocks: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnError(simErr)
			},
			expectedNames: nil,
			shouldErr:     true,
			errFragment:   "distinct query failed",
		},
		{
			title: "distinct scan error",
			setupMocks: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(sqlmock.NewRows([]string{"name", "extra"}).AddRow("X", "Y"))
			},
			expectedNames: nil,
			shouldErr:     true,
			errFragment:   "distinct row scan error",
		},
		{
			title: "distinct iteration error",
			setupMocks: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Z").
					RowError(0, simErr)
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			shouldErr:     true,
			errFragment:   "distinct iteration error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := db.New(mockDB)
			tc.setupMocks(mock)

			names, err := service.GetUniqueNames()

			if tc.shouldErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errFragment)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}