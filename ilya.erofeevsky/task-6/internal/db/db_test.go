package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ilya-Er0fick/task-6/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestDBService_GetNames(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func(sqlmock.Sqlmock)
		wantNames  []string
		wantErr    bool
		errContains string
	}{
		{
			name: "success with multiple rows",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob")
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{"Alice", "Bob"},
			wantErr:   false,
		},
		{
			name: "success with empty result",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{},
			wantErr:   false,
		},
		{
			name: "query error",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT name FROM users").WillReturnError(errors.New("db error"))
			},
			wantErr:     true,
			errContains: "query error",
		},
		{
			name: "scan error",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(123)
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantErr:     true,
			errContains: "scan error",
		},
		{
			name: "rows iteration error",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					RowError(0, errors.New("row error"))
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantErr:     true,
			errContains: "rows iteration error",
		},
		{
			name: "close error",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					CloseError(errors.New("close error"))
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantErr:     true,
			errContains: "close error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbConn, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer dbConn.Close()

			tt.setupMock(mock)
			service := db.New(dbConn)

			names, err := service.GetNames()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantNames, names)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func(sqlmock.Sqlmock)
		wantNames  []string
		wantErr    bool
		errContains string
	}{
		{
			name: "success with multiple rows",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob")
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{"Alice", "Bob"},
			wantErr:   false,
		},
		{
			name: "success with empty result",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{},
			wantErr:   false,
		},
		{
			name: "query error",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errors.New("db error"))
			},
			wantErr:     true,
			errContains: "query error",
		},
		{
			name: "scan error",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(123)
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantErr:     true,
			errContains: "scan error",
		},
		{
			name: "rows iteration error",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					RowError(0, errors.New("row error"))
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantErr:     true,
			errContains: "rows iteration error",
		},
		{
			name: "close error",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					CloseError(errors.New("close error"))
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantErr:     true,
			errContains: "close error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbConn, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer dbConn.Close()

			tt.setupMock(mock)
			service := db.New(dbConn)

			names, err := service.GetUniqueNames()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantNames, names)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
