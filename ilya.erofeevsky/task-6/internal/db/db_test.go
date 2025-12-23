package db_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ilya-Er0fick/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

var (
	errTestQuery = errors.New("query failed")
	errTestRows  = errors.New("rows error")
	errTestClose = errors.New("close error")
)

func TestDBService_AllScenarios(t *testing.T) {
	t.Parallel()

	methods := []string{"GetNames", "GetUniqueNames"}

	for _, method := range methods {
		method := method

		t.Run(method, func(t *testing.T) {
			t.Parallel()

			queryRegex := regexp.QuoteMeta("SELECT name FROM users")
			if method == "GetUniqueNames" {
				queryRegex = regexp.QuoteMeta("SELECT DISTINCT name FROM users")
			}

			tests := []struct {
				name    string
				setup   func(m sqlmock.Sqlmock)
				wantErr bool
			}{
				{
					name: "success",
					setup: func(m sqlmock.Sqlmock) {
						rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
						m.ExpectQuery(queryRegex).WillReturnRows(rows)
					},
					wantErr: false,
				},
				{
					name: "query_error",
					setup: func(m sqlmock.Sqlmock) {
						m.ExpectQuery(queryRegex).WillReturnError(errTestQuery)
					},
					wantErr: true,
				},
				{
					name: "scan_error",
					setup: func(m sqlmock.Sqlmock) {
						rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
						m.ExpectQuery(queryRegex).WillReturnRows(rows)
					},
					wantErr: true,
				},
				{
					name: "rows_iteration_error",
					setup: func(m sqlmock.Sqlmock) {
						rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errTestRows)
						m.ExpectQuery(queryRegex).WillReturnRows(rows)
					},
					wantErr: true,
				},
				{
					name: "close_error",
					setup: func(m sqlmock.Sqlmock) {
						rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").CloseError(errTestClose)
						m.ExpectQuery(queryRegex).WillReturnRows(rows)
					},
					wantErr: true,
				},
			}

			for _, tt := range tests {
				tt := tt
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					dbConn, mock, _ := sqlmock.New()
					defer dbConn.Close()

					service := db.New(dbConn)
					tt.setup(mock)

					var err error
					if method == "GetNames" {
						_, err = service.GetNames()
					} else {
						_, err = service.GetUniqueNames()
					}

					if tt.wantErr {
						require.Error(t, err)
					} else {
						require.NoError(t, err)
					}

					require.NoError(t, mock.ExpectationsWereMet())
				})
			}
		})
	}
}
