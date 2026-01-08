package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kostya.orlov/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

const (
	queryGetName   = "SELECT name FROM users"
	queryGetUnique = "SELECT DISTINCT name FROM users"
)

var ErrExpected = errors.New("error expected")

func TestGetNameSuccess(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}
	defer mockDB.Close()

	dbService := db.New(mockDB)

	testTable := [][]string{
		{"Peter", "Ivan", "Casey"},
		{"Jim89", "Sherlock76"},
		{"", ""},
	}

	for _, row := range testTable {
		mock.ExpectQuery(queryGetName).WillReturnRows(mockDBRows(row))

		names, err := dbService.GetNames()

		require.Equal(t, row, names)
		require.NoError(t, err)
	}
}

func TestGetNameDbQueryError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}
	defer mockDB.Close()

	dbService := db.New(mockDB)

	rows := mockDBRows(nil)

	mock.ExpectQuery(queryGetName).WillReturnRows(rows).WillReturnError(ErrExpected)

	names, err := dbService.GetNames()

	require.Nil(t, names)
	require.ErrorIs(t, err, ErrExpected)
	require.ErrorContains(t, err, "db query")
}

func TestGetNameScanError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}
	defer mockDB.Close()

	dbService := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow(nil)

	mock.ExpectQuery(queryGetName).WillReturnRows(rows)

	names, err := dbService.GetNames()

	require.Nil(t, names)
	require.ErrorContains(t, err, "rows scanning")
}

func TestGetNameRowsError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}
	defer mockDB.Close()

	dbService := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow("Peter")
	rows.RowError(0, ErrExpected)

	mock.ExpectQuery(queryGetName).WillReturnRows(rows)

	names, err := dbService.GetNames()

	require.Nil(t, names)
	require.ErrorContains(t, err, "rows error")
}

func TestGetUniqueNameSuccess(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}
	defer mockDB.Close()

	dbService := db.New(mockDB)

	testTable := [][]string{
		{"Peter", "Peter", "Peter"},
		{"Peter", "Peter", "Casey", "Casey", "Casey1"},
		{"", ""},
	}

	for _, row := range testTable {
		uniqueRow := uniqueRows(row)

		mock.ExpectQuery(queryGetUnique).WillReturnRows(mockDBRows(uniqueRow))

		names, err := dbService.GetUniqueNames()

		require.Equal(t, uniqueRow, names)
		require.NoError(t, err)
	}
}

func TestGetUniqueNameDbQueryError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}
	defer mockDB.Close()

	dbService := db.New(mockDB)

	rows := mockDBRows(nil)

	mock.ExpectQuery(queryGetUnique).WillReturnRows(rows).WillReturnError(ErrExpected)

	names, err := dbService.GetUniqueNames()

	require.Nil(t, names)
	require.ErrorIs(t, err, ErrExpected)
	require.ErrorContains(t, err, "db query")
}

func TestGetUniqueNameScanError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}
	defer mockDB.Close()

	dbService := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow(nil)

	mock.ExpectQuery(queryGetUnique).WillReturnRows(rows)

	names, err := dbService.GetUniqueNames()

	require.Nil(t, names)
	require.ErrorContains(t, err, "rows scanning")
}

func TestGetUniqueNameRowsError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}
	defer mockDB.Close()

	dbService := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow("Peter")
	rows.RowError(0, ErrExpected)

	mock.ExpectQuery(queryGetUnique).WillReturnRows(rows)

	names, err := dbService.GetUniqueNames()

	require.Nil(t, names)
	require.ErrorContains(t, err, "rows error")
}

func mockDBRows(names []string) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"name"})
	for _, name := range names {
		rows = rows.AddRow(name)
	}

	return rows
}

func uniqueRows(names []string) []string {
	find := func(value string, array []string) bool {
		for _, str := range array {
			if str == value {
				return true
			}
		}

		return false
	}

	var uniqueRows []string

	for _, name := range names {
		if !find(name, uniqueRows) {
			uniqueRows = append(uniqueRows, name)
		}
	}

	return uniqueRows
}