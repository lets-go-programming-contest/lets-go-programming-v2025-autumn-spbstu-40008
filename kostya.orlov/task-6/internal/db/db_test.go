package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/task-6/internal/db"
)

func TestDBService(t *testing.T) {
	errMock := errors.New("mock error")

	t.Run("GetNames", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			dbMock, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbMock.Close()

			rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
			mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

			service := db.New(dbMock)
			names, err := service.GetNames()

			assert.NoError(t, err)
			assert.Equal(t, []string{"Alice", "Bob"}, names)
		})

		t.Run("QueryError", func(t *testing.T) {
			dbMock, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbMock.Close()

			mock.ExpectQuery("SELECT name FROM users").WillReturnError(errMock)

			service := db.New(dbMock)
			_, err = service.GetNames()
			assert.Error(t, err)
		})

		t.Run("ScanError", func(t *testing.T) {
			dbMock, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbMock.Close()

			rows := sqlmock.NewRows([]string{"name", "extra"}).AddRow("Alice", "something")
			mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

			service := db.New(dbMock)
			_, err = service.GetNames()
			assert.Error(t, err)
		})

		t.Run("RowsIterationError", func(t *testing.T) {
			dbMock, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbMock.Close()

			rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errMock)
			mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

			service := db.New(dbMock)
			_, err = service.GetNames()
			assert.Error(t, err)
		})
	})

	t.Run("GetUniqueNames", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			dbMock, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbMock.Close()

			rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
			mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

			service := db.New(dbMock)
			names, err := service.GetUniqueNames()

			assert.NoError(t, err)
			assert.Equal(t, []string{"Alice"}, names)
		})

		t.Run("QueryError", func(t *testing.T) {
			dbMock, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbMock.Close()

			mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errMock)

			service := db.New(dbMock)
			_, err = service.GetUniqueNames()
			assert.Error(t, err)
		})

		t.Run("ScanError", func(t *testing.T) {
			dbMock, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbMock.Close()

			rows := sqlmock.NewRows([]string{"name", "extra"}).AddRow("Alice", "extra")
			mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

			service := db.New(dbMock)
			_, err = service.GetUniqueNames()
			assert.Error(t, err)
		})

		t.Run("RowsIterationError", func(t *testing.T) {
			dbMock, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbMock.Close()

			rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errMock)
			mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

			service := db.New(dbMock)
			_, err = service.GetUniqueNames()
			assert.Error(t, err)
		})
	})
}