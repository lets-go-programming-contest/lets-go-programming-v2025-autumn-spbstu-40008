package db_test

import (
	"errors"
	"testing"

	"example_mock/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestGetNames(t *testing.T) {
	type behavior func(m sqlmock.Sqlmock)

	testTable := []struct {
		name        string
		mockBehavior behavior
		expected    []string
		expectErr   bool
	}{
		{
			name: "Success",
			mockBehavior: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Ivan").AddRow("Gena228")
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expected:  []string{"Ivan", "Gena228"},
			expectErr: false,
		},
		{
			name: "Query Error",
			mockBehavior: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT name FROM users").WillReturnError(errors.New("db error"))
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name: "Scan Error",
			mockBehavior: func(m sqlmock.Sqlmock) {
				// Передаем nil в колонку типа string, что вызовет ошибку Scan
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name: "Rows Iteration Error",
			mockBehavior: func(m sqlmock.Sqlmock) {
				// Имитируем ошибку, возникающую в процессе итерации rows.Next() -> rows.Err()
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Ivan").
					RowError(0, errors.New("iteration error"))
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			tc.mockBehavior(mock)

			dbService := db.New(mockDB)
			names, err := dbService.GetNames()

			if tc.expectErr {
				require.Error(t, err)
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
