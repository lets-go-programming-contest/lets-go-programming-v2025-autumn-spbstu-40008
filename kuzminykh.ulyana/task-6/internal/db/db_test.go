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

func TestGetNames_ScanError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).AddRow(123)

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	names, err := service.GetNames()
	require.Error(t, err)
	require.Nil(t, names)
	require.Contains(t, err.Error(), "rows scanning")
}

func TestGetNames_RowsError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Ulya").
		RowError(0, errors.New("network failure"))

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	names, err := service.GetNames()
	require.Error(t, err)
	require.Nil(t, names)
	require.Contains(t, err.Error(), "rows error")
}

func TestGetNames_Empty(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()
	service := New(mockDB)
	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
	names, err := service.GetNames()
	require.NoError(t, err)
	require.Empty(t, names)
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

func TestGetUniqueNames_ScanError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).AddRow(456)
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.Error(t, err)
	require.Nil(t, names)
	require.Contains(t, err.Error(), "rows scanning")
}

func TestGetUniqueNames_RowsError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Ulya").
		RowError(0, errors.New("io error"))

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.Error(t, err)
	require.Nil(t, names)
	require.Contains(t, err.Error(), "rows error")
}

func TestGetUniqueNames_Empty(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()
	service := New(mockDB)
	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
	names, err := service.GetUniqueNames()
	require.NoError(t, err)
	require.Empty(t, names)
}
