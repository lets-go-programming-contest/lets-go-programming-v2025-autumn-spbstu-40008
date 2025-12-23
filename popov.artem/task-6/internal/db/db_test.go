package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/task-6/internal/db"
)

var mockErr = errors.New("mocked database error")

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		wantNames   []string
		wantErr     bool
		errContains string
	}{
		{
			name: "success",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{"Alice", "Bob"},
			wantErr:   false,
		},
		{
			name: "query error",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT name FROM users").WillReturnError(mockErr)
			},
			wantNames:   nil,
			wantErr:     true,
			errContains: "mocked database error",
		},
		{
			name: "scan error (wrong column count)",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name", "age"}).AddRow("Alice", 30)
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantNames:   nil,
			wantErr:     true,
			errContains: "sql: expected 1 destination arguments, not 2",
		},
		{
			name: "rows.Err() returns error after successful scan",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Charlie").RowError(0, mockErr)
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantNames:   nil,
			wantErr:     true,
			errContains: "mocked database error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := db.New(mockDB)
			tc.setupMock(mock)

			names, err := service.GetNames()

			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		wantNames   []string
		wantErr     bool
		errContains string
	}{
		{
			name: "success",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{"Alice"},
			wantErr:   false,
		},
		{
			name: "query error",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(mockErr)
			},
			wantNames:   nil,
			wantErr:     true,
			errContains: "mocked database error",
		},
		{
			name: "scan error",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name", "extra"}).AddRow("X", "Y")
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantNames:   nil,
			wantErr:     true,
			errContains: "sql: expected 1 destination arguments, not 2",
		},
		{
			name: "rows.Err() error",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Z").RowError(0, mockErr)
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantNames:   nil,
			wantErr:     true,
			errContains: "mocked database error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := db.New(mockDB)
			tc.setupMock(mock)

			names, err := service.GetUniqueNames()

			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
