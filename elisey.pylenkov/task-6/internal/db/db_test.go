package db_test

import (
	"errors"
	"testing"

	"task-6/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() {
			_ = dbConn.Close()
		}()

		svc := db.New(dbConn)

		dbRows := sqlmock.NewRows([]string{"name"}).AddRow("Ivan").AddRow("Gena")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(dbRows)

		res, err := svc.GetNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"Ivan", "Gena"}, res)
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() {
			_ = dbConn.Close()
		}()

		svc := db.New(dbConn)

		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errors.New("error"))

		res, err := svc.GetNames()
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() {
			_ = dbConn.Close()
		}()

		svc := db.New(dbConn)

		dbRows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(dbRows)

		res, err := svc.GetNames()
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() {
			_ = dbConn.Close()
		}()

		svc := db.New(dbConn)

		dbRows := sqlmock.NewRows([]string{"name"}).AddRow("User").RowError(0, errors.New("row error"))
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(dbRows)

		res, err := svc.GetNames()
		require.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() {
			_ = dbConn.Close()
		}()

		svc := db.New(dbConn)

		dbRows := sqlmock.NewRows([]string{"name"}).AddRow("Admin")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(dbRows)

		res, err := svc.GetUniqueNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"Admin"}, res)
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() {
			_ = dbConn.Close()
		}()

		svc := db.New(dbConn)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errors.New("error"))

		res, err := svc.GetUniqueNames()
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() {
			_ = dbConn.Close()
		}()

		svc := db.New(dbConn)

		dbRows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(dbRows)

		res, err := svc.GetUniqueNames()
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() {
			_ = dbConn.Close()
		}()

		svc := db.New(dbConn)

		dbRows := sqlmock.NewRows([]string{"name"}).AddRow("User").RowError(0, errors.New("row error"))
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(dbRows)

		res, err := svc.GetUniqueNames()
		require.Error(t, err)
		assert.Nil(t, res)
	})
}
