package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errDBQuery = errors.New("query execution failed")
	errRowScan = errors.New("row scan error")
)

func TestDataService_FetchAllNames(t *testing.T) {
	t.Parallel()

	t.Run("success_case", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		svc := DataService{DB: dbConn}

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Mikhail").AddRow("Dmitry")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		res, err := svc.FetchAllNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"Mikhail", "Dmitry"}, res)
	})

	t.Run("query_failure", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		svc := DataService{DB: dbConn}

		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errDBQuery)

		res, err := svc.FetchAllNames()
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("scan_failure", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		svc := DataService{DB: dbConn}

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		res, err := svc.FetchAllNames()
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "scan row error")
	})

	t.Run("rows_error", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		svc := DataService{DB: dbConn}

		rows := sqlmock.NewRows([]string{"name"}).AddRow("User").RowError(0, errRowScan)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		res, err := svc.FetchAllNames()
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "rows iteration error")
	})
}

func TestDataService_FetchDistinctNames(t *testing.T) {
	t.Parallel()

	t.Run("success_case", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		svc := DataService{DB: dbConn}

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Admin")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		res, err := svc.FetchDistinctNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"Admin"}, res)
	})

	t.Run("query_failure", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		svc := DataService{DB: dbConn}

		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errDBQuery)

		res, err := svc.FetchDistinctNames()
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("scan_failure", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		svc := DataService{DB: dbConn}

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		res, err := svc.FetchDistinctNames()
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("rows_error", func(t *testing.T) {
		t.Parallel()

		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		svc := DataService{DB: dbConn}

		rows := sqlmock.NewRows([]string{"name"}).AddRow("User").RowError(0, errRowScan)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		res, err := svc.FetchDistinctNames()
		require.Error(t, err)
		assert.Nil(t, res)
	})
}
