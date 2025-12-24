package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"rabbitdfs/task-6/internal/db"
)

var (
	errQueryFail = errors.New("fail")
	errRowFail   = errors.New("err")
)

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		svc := db.New(dbConn)
		dbRows := sqlmock.NewRows([]string{"name"}).AddRow("Mikhail").AddRow("Dmitry")

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(dbRows)

		res, err := svc.GetNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"Mikhail", "Dmitry"}, res)
	})

	t.Run("query_fail", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		svc := db.New(dbConn)

		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errQueryFail)

		res, err := svc.GetNames()

		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("scan_fail", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		svc := db.New(dbConn)
		dbRows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(dbRows)

		res, err := svc.GetNames()

		require.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "rows scanning")
	})

	t.Run("rows_error", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		svc := db.New(dbConn)
		dbRows := sqlmock.NewRows([]string{"name"}).AddRow("User").RowError(0, errRowFail)

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(dbRows)

		res, err := svc.GetNames()

		require.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "rows error")
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		svc := db.New(dbConn)
		dbRows := sqlmock.NewRows([]string{"name"}).AddRow("Admin")

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(dbRows)

		res, err := svc.GetUniqueNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"Admin"}, res)
	})

	t.Run("query_fail", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		svc := db.New(dbConn)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errQueryFail)

		res, err := svc.GetUniqueNames()

		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("scan_fail", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		svc := db.New(dbConn)
		dbRows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(dbRows)

		res, err := svc.GetUniqueNames()

		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("rows_error", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		svc := db.New(dbConn)
		dbRows := sqlmock.NewRows([]string{"name"}).AddRow("User").RowError(0, errRowFail)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(dbRows)

		res, err := svc.GetUniqueNames()

		require.Error(t, err)
		assert.Nil(t, res)
	})
}
