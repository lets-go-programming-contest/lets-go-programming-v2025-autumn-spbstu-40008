package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errStub = errors.New("stub error")

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(mock sqlmock.Sqlmock)
		want      []string
		wantError bool
		errText   string
	}{
		{
			name: "Success",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Ivan").AddRow("Petr")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			want:      []string{"Ivan", "Petr"},
			wantError: false,
		},
		{
			name: "Query Error",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").WillReturnError(errStub)
			},
			want:      nil,
			wantError: true,
			errText:   "execution error",
		},
		{
			name: "Scan Error",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			want:      nil,
			wantError: true,
			errText:   "scan error",
		},
		{
			name: "Iteration Error",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Ivan").RowError(0, errStub)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			want:      nil,
			wantError: true,
			errText:   "iteration error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := New(mockDB)
			tc.setup(mock)

			got, err := service.GetNames()

			if tc.wantError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errText)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(mock sqlmock.Sqlmock)
		want      []string
		wantError bool
		errText   string
	}{
		{
			name: "Success",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("UniqueName")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			want:      []string{"UniqueName"},
			wantError: false,
		},
		{
			name: "Query Error",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errStub)
			},
			want:      nil,
			wantError: true,
			errText:   "execution error",
		},
		{
			name: "Scan Error",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name", "extra"}).AddRow("Name", "ExtraVal")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			want:      nil,
			wantError: true,
			errText:   "scan error",
		},
		{
			name: "Iteration Error",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Name").RowError(0, errStub)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			want:      nil,
			wantError: true,
			errText:   "iteration error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := New(mockDB)
			tc.setup(mock)

			got, err := service.GetUniqueNames()

			if tc.wantError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errText)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}