package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	internalDb "task-6/internal/db"
)

var (
	errConnection = errors.New("connection failed")
	errRow        = errors.New("row failure")
	errDBDead     = errors.New("db dead")
)

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	query := "SELECT name FROM users"

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		service := internalDb.New(db)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery(query).WillReturnRows(rows)

		names, err := service.GetNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Bob"}, names)
	})

	t.Run("success empty", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		service := internalDb.New(db)

		rows := sqlmock.NewRows([]string{"name"})

		mock.ExpectQuery(query).WillReturnRows(rows)

		names, err := service.GetNames()

		require.NoError(t, err)
		assert.Empty(t, names)
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		service := internalDb.New(db)

		mock.ExpectQuery(query).WillReturnError(errConnection)

		names, err := service.GetNames()

		require.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "db query")
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		service := internalDb.New(db)

		rows := sqlmock.NewRows([]string{"col1", "col2"}).AddRow("val1", "val2")

		mock.ExpectQuery(query).WillReturnRows(rows)

		names, err := service.GetNames()

		require.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "rows scanning")
	})

	t.Run("rows iteration error", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		service := internalDb.New(db)

		rows := sqlmock.NewRows([]string{"name"}).
			RowError(0, errRow)

		mock.ExpectQuery(query).WillReturnRows(rows)

		names, err := service.GetNames()

		require.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "rows error")
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	query := "SELECT DISTINCT name FROM users"

	// 1. Успешное получение данных
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		service := internalDb.New(db)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Charlie").
			AddRow("Dave")

		mock.ExpectQuery(query).WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"Charlie", "Dave"}, names)
	})

	t.Run("success empty", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		service := internalDb.New(db)

		rows := sqlmock.NewRows([]string{"name"})

		mock.ExpectQuery(query).WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		require.NoError(t, err)
		assert.Empty(t, names)
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		service := internalDb.New(db)

		mock.ExpectQuery(query).WillReturnError(errDBDead)

		names, err := service.GetUniqueNames()

		require.Error(t, err)
		assert.Nil(t, names)
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		service := internalDb.New(db)

		rows := sqlmock.NewRows([]string{"col1", "col2"}).AddRow("val1", "val2")

		mock.ExpectQuery(query).WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		require.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "rows scanning")
	})

	t.Run("rows iteration error", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		service := internalDb.New(db)

		rows := sqlmock.NewRows([]string{"name"}).
			RowError(0, errRow)

		mock.ExpectQuery(query).WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		require.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "rows error")
	})
}
