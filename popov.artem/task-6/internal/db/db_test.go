package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/task-6/internal/db"
)

var mockErr = errors.New("simulated database error")

func TestFetchUserNames(t *testing.T) {
	t.Parallel()

	cases := []struct {
		desc      string
		setupMock func(sqlmock.Sqlmock)
		wantNames []string
		expectErr bool
		errSubstr string
	}{
		{
			desc: "returns correct names",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT name FROM users").
					WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Eve").AddRow("John"))
			},
			wantNames: []string{"Eve", "John"},
			expectErr: false,
		},
		{
			desc: "query fails",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT name FROM users").
					WillReturnError(mockErr)
			},
			wantNames: nil,
			expectErr: true,
			errSubstr: "failed to execute query",
		},
		{
			desc: "scan fails due to column mismatch",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT name FROM users").
					WillReturnRows(sqlmock.NewRows([]string{"name", "age"}).AddRow("Alice", 30))
			},
			wantNames: nil,
			expectErr: true,
			errSubstr: "failed to scan row",
		},
		{
			desc: "row iteration error",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					RowError(0, mockErr)
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantNames: nil,
			expectErr: true,
			errSubstr: "error during row iteration",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			dbMock, sqlMock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbMock.Close()

			service := db.NewService(dbMock)
			tc.setupMock(sqlMock)

			names, err := service.FetchUserNames()

			if tc.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errSubstr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantNames, names)
			}

			require.NoError(t, sqlMock.ExpectationsWereMet())
		})
	}
}

func TestFetchDistinctUserNames(t *testing.T) {
	t.Parallel()

	cases := []struct {
		desc      string
		setupMock func(sqlmock.Sqlmock)
		wantNames []string
		expectErr bool
		errSubstr string
	}{
		{
			desc: "distinct names returned",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Bob"))
			},
			wantNames: []string{"Bob"},
			expectErr: false,
		},
		{
			desc: "distinct query fails",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnError(mockErr)
			},
			wantNames: nil,
			expectErr: true,
			errSubstr: "failed to execute distinct query",
		},
		{
			desc: "scan error on distinct",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(sqlmock.NewRows([]string{"name", "extra"}).AddRow("X", "Y"))
			},
			wantNames: nil,
			expectErr: true,
			errSubstr: "scan failed for distinct name",
		},
		{
			desc: "row error in distinct",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Z").
					RowError(0, mockErr)
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantNames: nil,
			expectErr: true,
			errSubstr: "row error in distinct query",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			dbMock, sqlMock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbMock.Close()

			service := db.NewService(dbMock)
			tc.setupMock(sqlMock)

			names, err := service.FetchDistinctUserNames()

			if tc.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errSubstr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantNames, names)
			}

			require.NoError(t, sqlMock.ExpectationsWereMet())
		})
	}
}
