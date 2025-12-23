package db_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ilya-Er0fick/task-6/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestDBService_FullCoverage(t *testing.T) {
	t.Parallel()

	queryNames := regexp.QuoteMeta("SELECT name FROM users")
	queryUnique := regexp.QuoteMeta("SELECT DISTINCT name FROM users")

	tests := []struct {
		name       string
		method     func(s db.DBService) ([]string, error)
		query      string
		setupMock  func(m sqlmock.Sqlmock, q string)
		want       []string
		wantErrMsg string
	}{

		{
			name:   "GetNames success",
			method: func(s db.DBService) ([]string, error) { return s.GetNames() },
			query:  queryNames,
			setupMock: func(m sqlmock.Sqlmock, q string) {
				m.ExpectQuery(q).WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("User1"))
			},
			want: []string{"User1"},
		},
		{
			name:   "GetNames query error",
			method: func(s db.DBService) ([]string, error) { return s.GetNames() },
			query:  queryNames,
			setupMock: func(m sqlmock.Sqlmock, q string) {
				m.ExpectQuery(q).WillReturnError(errors.New("fail"))
			},
			wantErrMsg: "query error",
		},
		{
			name:   "GetNames scan error",
			method: func(s db.DBService) ([]string, error) { return s.GetNames() },
			query:  queryNames,
			setupMock: func(m sqlmock.Sqlmock, q string) {
				m.ExpectQuery(q).WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
			},
			wantErrMsg: "scan error",
		},
		{
			name:   "GetNames iteration error",
			method: func(s db.DBService) ([]string, error) { return s.GetNames() },
			query:  queryNames,
			setupMock: func(m sqlmock.Sqlmock, q string) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("U1").RowError(0, errors.New("fail"))
				m.ExpectQuery(q).WillReturnRows(rows)
			},
			wantErrMsg: "rows iteration error",
		},
		{
			name:   "GetNames close error",
			method: func(s db.DBService) ([]string, error) { return s.GetNames() },
			query:  queryNames,
			setupMock: func(m sqlmock.Sqlmock, q string) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("U1").CloseError(errors.New("fail"))
				m.ExpectQuery(q).WillReturnRows(rows)
			},
			wantErrMsg: "close error",
		},

		{
			name:   "GetUniqueNames success",
			method: func(s db.DBService) ([]string, error) { return s.GetUniqueNames() },
			query:  queryUnique,
			setupMock: func(m sqlmock.Sqlmock, q string) {
				m.ExpectQuery(q).WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("U1"))
			},
			want: []string{"U1"},
		},
		{
			name:   "GetUniqueNames query error",
			method: func(s db.DBService) ([]string, error) { return s.GetUniqueNames() },
			query:  queryUnique,
			setupMock: func(m sqlmock.Sqlmock, q string) {
				m.ExpectQuery(q).WillReturnError(errors.New("fail"))
			},
			wantErrMsg: "query error",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dbConn, mock, _ := sqlmock.New()
			defer dbConn.Close()

			service := db.New(dbConn)
			tt.setupMock(mock, tt.query)

			res, err := tt.method(service)

			if tt.wantErrMsg != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, res)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
