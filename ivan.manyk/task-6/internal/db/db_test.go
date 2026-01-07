package db_test

import (
	"errors"
	"testing"

	database "task-6/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errConnection = errors.New("connection error")
	errRow        = errors.New("row error")
	errDatabase   = errors.New("database error")
)

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupMock    func(mock sqlmock.Sqlmock)
		expectedErr  bool
		expectedRows []string
	}{
		{
			name: "successful get names",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Ivan").
					AddRow("Petr").
					AddRow("Artem")
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rows)
			},
			expectedErr:  false,
			expectedRows: []string{"Ivan", "Petr", "Artem"},
		},
		{
			name: "query error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnError(errConnection)
			},
			expectedErr:  true,
			expectedRows: nil,
		},
		{
			name: "scan error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Ivan").
					AddRow(nil).
					AddRow("Artem")
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rows)
			},
			expectedErr:  true,
			expectedRows: nil,
		},
		{
			name: "rows error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Ivan").
					RowError(0, errRow)
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rows)
			},
			expectedErr:  true,
			expectedRows: nil,
		},
		{
			name: "empty result",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rows)
			},
			expectedErr:  false,
			expectedRows: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.setupMock(mock)

			service := database.New(db)

			result, err := service.GetNames()

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				if tt.expectedRows == nil {
					assert.Nil(t, result)
				} else {
					assert.Equal(t, tt.expectedRows, result)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupMock    func(mock sqlmock.Sqlmock)
		expectedErr  bool
		expectedRows []string
	}{
		{
			name: "successful get unique names",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Ivan").
					AddRow("Ivan").
					AddRow("Petr").
					AddRow("Artem").
					AddRow("Petr")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(rows)
			},
			expectedErr:  false,
			expectedRows: []string{"Ivan", "Ivan", "Petr", "Artem", "Petr"},
		},
		{
			name: "query error for unique names",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnError(errDatabase)
			},
			expectedErr:  true,
			expectedRows: nil,
		},
		{
			name: "scan error for unique names",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Ivan").
					AddRow(nil).
					AddRow("Artem")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(rows)
			},
			expectedErr:  true,
			expectedRows: nil,
		},
		{
			name: "rows error for unique names",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Ivan").
					RowError(0, errRow)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(rows)
			},
			expectedErr:  true,
			expectedRows: nil,
		},
		{
			name: "empty result for unique names",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(rows)
			},
			expectedErr:  false,
			expectedRows: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.setupMock(mock)

			service := database.New(db)

			result, err := service.GetUniqueNames()

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				if tt.expectedRows == nil {
					assert.Nil(t, result)
				} else {
					assert.Equal(t, tt.expectedRows, result)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
