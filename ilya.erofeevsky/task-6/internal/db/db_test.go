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
	errQuery = errors.New("query failed")
	errScan  = errors.New("scan failed")
)

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(m sqlmock.Sqlmock)
		want    []string
		wantErr bool
	}{
		{
			name: "success",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
				m.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)
			},
			want: []string{"Alice", "Bob"},
		},
		{
			name: "query error",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnError(errQuery)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbConn, mock, _ := sqlmock.New()
			defer dbConn.Close()
			
			svc := db.New(dbConn)
			tt.setup(mock)

			got, err := svc.GetNames()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
