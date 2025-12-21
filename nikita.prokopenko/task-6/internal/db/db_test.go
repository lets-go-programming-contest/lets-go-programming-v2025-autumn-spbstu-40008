package db

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockDB struct {
	queryFunc func(query string, args ...any) (*sql.Rows, error)
}

func (m *mockDB) Query(query string, args ...any) (*sql.Rows, error) {
	return m.queryFunc(query, args...)
}

func TestDataHandler_GetNames(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	handler := New(db)

	t.Run("success with multiple rows", func(t *testing.T) {
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
		mock.ExpectQuery("SELECT name FROM users").
			WillReturnError(errors.New("database error"))

		result, err := handler.GetNames()
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "database query failed")
	})

	t.Run("no records found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"})
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		result, err := handler.GetNames()
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "no records found")
	})

	t.Run("error on rows scan", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
		rows = rows.AddRow(nil)
		rows = rows.AddRow("Charlie")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		result, err := handler.GetNames()
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "row processing error")
	})

	t.Run("error on rows iteration", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob").
			RowError(1, errors.New("row error"))
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		result, err := handler.GetNames()
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "row processing error")
	})
}

func TestDataHandler_GetUniqueNames(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	handler := New(db)

	t.Run("success with duplicates", func(t *testing.T) {
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
		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnError(errors.New("database error"))

		result, err := handler.GetUniqueNames()
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "database query failed")
	})

	t.Run("no records found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"})
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		result, err := handler.GetUniqueNames()
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "no records found")
	})

	t.Run("error on rows scan", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow(nil).
			AddRow("Bob")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		result, err := handler.GetUniqueNames()
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "row processing error")
	})

	t.Run("error on rows iteration", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob").
			RowError(1, errors.New("row error"))
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		result, err := handler.GetUniqueNames()
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "row processing error")
	})

	t.Run("success with single empty name", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		result, err := handler.GetUniqueNames()
		require.NoError(t, err)
		assert.Equal(t, []string{""}, result)
	})
}

func TestDataHandler_New(t *testing.T) {
	t.Parallel()

	mockDB := &mockDB{
		queryFunc: func(query string, args ...any) (*sql.Rows, error) {
			return nil, nil
		},
	}
	handler := New(mockDB)
	assert.NotNil(t, handler)
	assert.Equal(t, mockDB, handler.DB)
}

func TestDataHandler_WithEmptyResult(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	handler := New(db)

	rows := sqlmock.NewRows([]string{"name"}).AddRow("test")
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	result, err := handler.GetNames()
	require.NoError(t, err)
	assert.Equal(t, []string{"test"}, result)
}

func TestDataHandler_GetUniqueNamesWithEmptyString(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	handler := New(db)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("").
		AddRow("").
		AddRow("")
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	result, err := handler.GetUniqueNames()
	require.NoError(t, err)
	assert.Equal(t, []string{""}, result)
}