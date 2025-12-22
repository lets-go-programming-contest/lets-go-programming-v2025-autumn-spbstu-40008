package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mordw1n/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

var (
	errDatabase   = errors.New("database error")
	errConstraint = errors.New("constraint violation")
)

func TestGetNames(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	t.Cleanup(func() { mockDB.Close() })

	dbService := db.New(mockDB)

	tests := []struct {
		name      string
		setupMock func()
		wantNames []string
		wantError bool
		errorMsg  string
	}{
		{
			name: "successful query with multiple rows",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Mikhail").
					AddRow("Phillip").
					AddRow("Alex")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{"Mikhail", "Phillip", "Alex"},
			wantError: false,
		},
		{
			name: "empty result",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{},
			wantError: false,
		},
		{
			name: "query error",
			setupMock: func() {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnError(errDatabase)
			},
			wantNames: nil,
			wantError: true,
			errorMsg:  "db query:",
		},
		{
			name: "scan error in one row",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Mikhail").
					AddRow(nil).
					AddRow("Alex")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantNames: nil,
			wantError: true,
			errorMsg:  "rows scanning:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMock()

			names, err := dbService.GetNames()

			if tt.wantError {
				require.Error(t, err)

				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	t.Cleanup(func() { mockDB.Close() })

	dbService := db.New(mockDB)

	tests := []struct {
		name      string
		setupMock func()
		wantNames []string
		wantError bool
		errorMsg  string
	}{
		{
			name: "successful query with multiple unique rows",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Mikhail").
					AddRow("Phillip").
					AddRow("Alex")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{"Mikhail", "Phillip", "Alex"},
			wantError: false,
		},
		{
			name: "empty result",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{},
			wantError: false,
		},
		{
			name: "query error",
			setupMock: func() {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnError(errDatabase)
			},
			wantNames: nil,
			wantError: true,
			errorMsg:  "db query:",
		},
		{
			name: "scan error in one row",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Mikhail").
					AddRow(nil).
					AddRow("Alex")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantNames: nil,
			wantError: true,
			errorMsg:  "rows scanning:",
		},
		{
			name: "duplicate names in source (should return unique)",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Mikhail").
					AddRow("Mikhail").
					AddRow("Alex").
					AddRow("Alex").
					AddRow("Phillip")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{"Mikhail", "Mikhail", "Alex", "Alex", "Phillip"},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMock()

			names, err := dbService.GetUniqueNames()

			if tt.wantError {
				require.Error(t, err)

				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)

	t.Cleanup(func() { mockDB.Close() })

	service := db.New(mockDB)
	require.NotNil(t, service)
}

func TestRowsErrorHandling(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	t.Cleanup(func() { mockDB.Close() })

	dbService := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Mikhail").
		AddRow("Alex").
		RowError(1, errors.New("row error"))

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	names, err := dbService.GetNames()
	require.Error(t, err)
	require.Nil(t, names)
	require.Contains(t, err.Error(), "rows error:")
	require.NoError(t, mock.ExpectationsWereMet())
}
