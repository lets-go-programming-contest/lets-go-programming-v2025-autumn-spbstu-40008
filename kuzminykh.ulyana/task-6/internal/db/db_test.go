package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

type rowTestDb struct {
	names       []string
	errExpected error
}

var testTable = []rowTestDb{
	{
		names: []string{"Ulya", "Keks"},
	},
	{
		errExpected: errors.New("database query error"),
	},
}

func TestGetNames(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := New(mockDB)

	for i, row := range testTable {
		rows := sqlmock.NewRows([]string{"name"})
		for _, name := range row.names {
			rows = rows.AddRow(name)
		}

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(rows).
			WillReturnError(row.errExpected)

		names, err := service.GetNames()

		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected, "row: %d", i)
			require.Nil(t, names, "row: %d", i)
			continue
		}

		require.NoError(t, err, "row: %d", i)
		require.Equal(t, row.names, names, "row: %d", i)
	}
}

func TestGetUniqueNames(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := New(mockDB)

	for i, row := range testTable {
		rows := sqlmock.NewRows([]string{"name"})
		for _, name := range row.names {
			rows = rows.AddRow(name)
		}

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(rows).
			WillReturnError(row.errExpected)

		names, err := service.GetUniqueNames()

		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected, "row: %d", i)
			require.Nil(t, names, "row: %d", i)
			continue
		}

		require.NoError(t, err, "row: %d", i)
		require.Equal(t, row.names, names, "row: %d", i)
	}
}

func mockDbRows(names []string) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"name"})
	for _, name := range names {
		rows = rows.AddRow(name)
	}
	return rows
}
