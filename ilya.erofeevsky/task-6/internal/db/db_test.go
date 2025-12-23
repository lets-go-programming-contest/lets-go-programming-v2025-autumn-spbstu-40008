package db_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ilya-Er0fick/task-6/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()
	query := regexp.QuoteMeta("SELECT name FROM users")

	tests := []struct {
		name    string
		setup   func(sqlmock.Sqlmock)
		want    []string
		wantErr bool
	}{
		{
			name: "success",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("User1").AddRow("User2")
				m.ExpectQuery(query).WillReturnRows(rows)
			},
			want: []string{"User1", "User2"},
		},
		{
			name: "query_error",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(query).WillReturnError(errors.New("db fail"))
			},
			wantErr: true,
		},
		{
			name: "scan_error",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil) // trigger scan error
				m.ExpectQuery(query).WillReturnRows(rows)
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

			res, err := service.GetNames()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, res)
			}
		})
	}
}
