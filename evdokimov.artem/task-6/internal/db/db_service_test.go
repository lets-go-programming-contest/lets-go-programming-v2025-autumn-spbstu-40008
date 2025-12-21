package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	dbpkg "github.com/LAffey26/task-6/internal/db"
)

var (
	errConnection = errors.New("connection error")
	errRow        = errors.New("row error")
)

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name        string
		mockSetup   func(sqlmock.Sqlmock)
		want        []string
		wantErrPart string
	}

	testCases := []testCase{
		{
			name: "ok",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Artem").
					AddRow("Bob")

				mock.ExpectQuery("^SELECT name FROM users$").
					WillReturnRows(rows)
			},
			want: []string{"Artem", "Bob"},
		},
		{
			name: "ok empty",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})

				mock.ExpectQuery("^SELECT name FROM users$").
					WillReturnRows(rows)
			},
			want: nil,
		},
		{
			name: "query error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT name FROM users$").
					WillReturnError(errConnection)
			},
			wantErrPart: "db query:",
		},
		{
			name: "scan error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow(nil)

				mock.ExpectQuery("^SELECT name FROM users$").
					WillReturnRows(rows)
			},
			wantErrPart: "rows scanning:",
		},
		{
			name: "rows error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Artem").
					AddRow("Bob").
					RowError(1, errRow)

				mock.ExpectQuery("^SELECT name FROM users$").
					WillReturnRows(rows)
			},
			wantErrPart: "rows error:",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			t.Cleanup(func() { _ = mockDB.Close() })

			service := dbpkg.New(mockDB)

			tc.mockSetup(mock)

			got, gotErr := service.GetNames()

			if tc.wantErrPart != "" {
				require.Error(t, gotErr)
				require.Contains(t, gotErr.Error(), tc.wantErrPart)
				require.Nil(t, got)
			} else {
				require.NoError(t, gotErr)
				require.Equal(t, tc.want, got)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name        string
		mockSetup   func(sqlmock.Sqlmock)
		want        []string
		wantErrPart string
	}

	testCases := []testCase{
		{
			name: "ok",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Artem").
					AddRow("Bob")

				mock.ExpectQuery("^SELECT DISTINCT name FROM users$").
					WillReturnRows(rows)
			},
			want: []string{"Artem", "Bob"},
		},
		{
			name: "query error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT DISTINCT name FROM users$").
					WillReturnError(sql.ErrConnDone)
			},
			wantErrPart: "db query:",
		},
		{
			name: "scan error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow(nil)

				mock.ExpectQuery("^SELECT DISTINCT name FROM users$").
					WillReturnRows(rows)
			},
			wantErrPart: "rows scanning:",
		},
		{
			name: "rows error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Artem").
					AddRow("Bob").
					RowError(1, errRow)

				mock.ExpectQuery("^SELECT DISTINCT name FROM users$").
					WillReturnRows(rows)
			},
			wantErrPart: "rows error:",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			t.Cleanup(func() { _ = mockDB.Close() })

			service := dbpkg.New(mockDB)

			tc.mockSetup(mock)

			got, gotErr := service.GetUniqueNames()

			if tc.wantErrPart != "" {
				require.Error(t, gotErr)
				require.Contains(t, gotErr.Error(), tc.wantErrPart)
				require.Nil(t, got)
			} else {
				require.NoError(t, gotErr)
				require.Equal(t, tc.want, got)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = mockDB.Close() })

	service := dbpkg.New(mockDB)

	require.Equal(t, mockDB, service.DB)
}
