package db

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataHandler_GetNames(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	handler := New(db)

	t.Run("success with multiple rows", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob").
			AddRow("Charlie")

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		result, err := handler.GetNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Bob", "Charlie"}, result)
	})

	t.Run("error on query execution", func(t *testing.T) {
		t.Parallel()

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnError(errors.New("database error"))

		result, err := handler.GetNames()
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "database query failed")
	})

	t.Run("no records found", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"})
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		result, err := handler.GetNames()
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "no records found")
	})
}

func TestDataHandler_GetUniqueNames(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	handler := New(db)

	t.Run("success with duplicates", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob").
			AddRow("Alice").
			AddRow("Charlie").
			AddRow("Bob")

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		result, err := handler.GetUniqueNames()
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{"Alice", "Bob", "Charlie"}, result)
	})

	t.Run("error on query execution", func(t *testing.T) {
		t.Parallel()

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnError(errors.New("database error"))

		result, err := handler.GetUniqueNames()
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "database query failed")
	})

	t.Run("no records found", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"})
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		result, err := handler.GetUniqueNames()
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "no records found")
	})
}