package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mordw1n/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

var (
	errRow      = errors.New("row error")
	errDatabase = errors.New("database error")
)

func TestGetNames(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	t.Cleanup(func() { mockDB.Close() })

	dbService := db.New(mockDB)

	t.Run("successful query with multiple rows", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Mikhail").
			AddRow("Phillip").
			AddRow("Alex")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := dbService.GetNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Mikhail", "Phillip", "Alex"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"})
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := dbService.GetNames()
		require.NoError(t, err)
		require.Equal(t, []string{}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnError(errDatabase)

		names, err := dbService.GetNames()
		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "db query:")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error with nil value", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Mikhail").
			AddRow(nil).
			AddRow("Alex")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := dbService.GetNames()
		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "rows scanning:")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error with wrong type", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow(123)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := dbService.GetNames()
		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "rows scanning:")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error during iteration", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Mikhail").
			AddRow("Alex").
			RowError(1, errRow)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := dbService.GetNames()
		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "rows error:")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error before iteration", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"}).
			RowError(0, errRow)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := dbService.GetNames()
		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "rows error:")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("connection closed before query", func(t *testing.T) {
		t.Parallel()

		mockDB2, err := sqlmock.New()
		require.NoError(t, err)
		mockDB2.Close()

		dbService2 := db.New(mockDB2)

		names, err := dbService2.GetNames()
		require.Error(t, err)
		require.Nil(t, names)
	})
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	t.Cleanup(func() { mockDB.Close() })

	dbService := db.New(mockDB)

	t.Run("successful query with multiple unique rows", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Mikhail").
			AddRow("Phillip").
			AddRow("Alex")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := dbService.GetUniqueNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Mikhail", "Phillip", "Alex"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"})
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := dbService.GetUniqueNames()
		require.NoError(t, err)
		require.Equal(t, []string{}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnError(errDatabase)

		names, err := dbService.GetUniqueNames()
		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "db query:")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error with nil value", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Mikhail").
			AddRow(nil).
			AddRow("Alex")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := dbService.GetUniqueNames()
		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "rows scanning:")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error with wrong type", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow(123)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := dbService.GetUniqueNames()
		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "rows scanning:")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error during iteration", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Mikhail").
			AddRow("Alex").
			RowError(1, errRow)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := dbService.GetUniqueNames()
		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "rows error:")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error before iteration", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"}).
			RowError(0, errRow)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := dbService.GetUniqueNames()
		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "rows error:")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate names in source", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Mikhail").
			AddRow("Alex").
			AddRow("Mikhail").
			AddRow("Phillip")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := dbService.GetUniqueNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Mikhail", "Alex", "Mikhail", "Phillip"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("connection closed before query", func(t *testing.T) {
		t.Parallel()

		mockDB2, err := sqlmock.New()
		require.NoError(t, err)
		mockDB2.Close()

		dbService2 := db.New(mockDB2)

		names, err := dbService2.GetUniqueNames()
		require.Error(t, err)
		require.Nil(t, names)
	})
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)

	t.Cleanup(func() { mockDB.Close() })

	service := db.New(mockDB)
	require.NotNil(t, service)
}

func TestMultipleCalls(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	t.Cleanup(func() { mockDB.Close() })

	dbService := db.New(mockDB)

	rows1 := sqlmock.NewRows([]string{"name"}).
		AddRow("Mikhail").
		AddRow("Phillip")
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows1)

	names1, err := dbService.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"Mikhail", "Phillip"}, names1)

	rows2 := sqlmock.NewRows([]string{"name"}).
		AddRow("Mikhail").
		AddRow("Alex")
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows2)

	names2, err := dbService.GetUniqueNames()
	require.NoError(t, err)
	require.Equal(t, []string{"Mikhail", "Alex"}, names2)

	require.NoError(t, mock.ExpectationsWereMet())
}
