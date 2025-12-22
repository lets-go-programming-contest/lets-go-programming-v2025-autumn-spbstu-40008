package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mordw1n/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

var (
	errDatabase = errors.New("database error")
	errRow      = errors.New("row error")
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

	t.Run("scan error in one row", func(t *testing.T) {
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

	t.Run("rows error after iteration", func(t *testing.T) {
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

	t.Run("scan error type mismatch", func(t *testing.T) {
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

	t.Run("scan error in one row", func(t *testing.T) {
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

	t.Run("rows error after iteration", func(t *testing.T) {
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

	t.Run("duplicate names in source (should return all with DISTINCT)", func(t *testing.T) {
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

	t.Run("scan error type mismatch", func(t *testing.T) {
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
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)

	t.Cleanup(func() { mockDB.Close() })

	service := db.New(mockDB)
	require.NotNil(t, service)
}
