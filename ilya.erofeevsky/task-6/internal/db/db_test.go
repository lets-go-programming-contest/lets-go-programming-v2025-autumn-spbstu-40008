package db_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ilya-Er0fick/task-6/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestDBService(t *testing.T) {
	t.Parallel()

	commonTests := []struct {
		name       string
		method     func(db.DBService) ([]string, error)
		query      string
		setupMock  func(m sqlmock.Sqlmock, q string)
		want       []string
		wantErrMsg string
	}{
		{
			name:   "GetNames success",
			method: func(s db.DBService) ([]string, error) { return s.GetNames() },
			query:  regexp.QuoteMeta("SELECT name FROM users"),
			setupMock: func(m sqlmock.Sqlmock, q string) {
				m.ExpectQuery(q).WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("User1"))
			},
			want: []string{"User1"},
		},
		{
			name:   "GetUniqueNames success",
			method: func(s db.DBService) ([]string, error) { return s.GetUniqueNames() },
			query:  regexp.QuoteMeta("SELECT DISTINCT name FROM users"),
			setupMock: func(m sqlmock.Sqlmock, q string) {
				m.ExpectQuery(q).WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Admin"))
			},
			want: []string{"Admin"},
		},
		{
			name:   "Query error",
			method: func(s db.DBService) ([]string, error) { return s.GetNames() },
			query:  regexp.QuoteMeta("SELECT name FROM users"),
			setupMock: func(m sqlmock.Sqlmock, q string) {
				m.ExpectQuery(q).WillReturnError(errors.New("fail"))
			},
			wantErrMsg: "query error",
		},
		{
			name:   "Scan error",
			method: func(s db.DBService) ([]string, error) { return s.GetNames() },
			query:  regexp.QuoteMeta("SELECT name FROM users"),
			setupMock: func(m sqlmock.Sqlmock, q string) {
				m.ExpectQuery(q).WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
			},
			wantErrMsg: "scan error",
		},
		{
			name:   "Rows iteration error",
			method: func(s db.DBService) ([]string, error) { return s.GetNames() },
			query:  regexp.QuoteMeta("SELECT name FROM users"),
			setupMock: func(m sqlmock.Sqlmock, q string) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("U1").RowError(0, errors.New("iter fail"))
				m.ExpectQuery(q).WillReturnRows(rows)
			},
			wantErrMsg: "rows iteration error",
		},
		{
			name:   "Close error",
			method: func(s db.DBService) ([]string, error) { return s.GetNames() },
			query:  regexp.QuoteMeta("SELECT name FROM users"),
			setupMock: func(m sqlmock.Sqlmock, q string) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("U1").CloseError(errors.New("close fail"))
				m.ExpectQuery(q).WillReturnRows(rows)
			},
			wantErrMsg: "close error",
		},
	}

	for _, tt := range commonTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sqlDB, mock, _ := sqlmock.New()
			defer sqlDB.Close()
			
			service := db.New(sqlDB)
			tt.setupMock(mock, tt.query)

			res, err := tt.method(service)

			if tt.wantErrMsg != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, res)
			}
		})
	}
}
