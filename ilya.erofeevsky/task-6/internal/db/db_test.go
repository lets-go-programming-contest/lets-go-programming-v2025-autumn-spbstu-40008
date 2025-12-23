package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ilya-Er0fick/task-6/internal/db"
	"github.com/stretchr/testify/assert"
)


func TestGetNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		service := db.New(dbConn)
		names, err := service.GetNames()

		assert.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Bob"}, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query_error", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errors.New("db error"))

		service := db.New(dbConn)
		_, err := service.GetNames()

		assert.Error(t, err)
	})

	t.Run("scan_error", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		service := db.New(dbConn)
		_, err := service.GetNames()

		assert.Error(t, err)
	})

	t.Run("rows_err_loop", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errors.New("loop error"))
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		service := db.New(dbConn)
		_, err := service.GetNames()

		assert.Error(t, err)
	})

	t.Run("close_error", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").CloseError(errors.New("close error"))
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		service := db.New(dbConn)
		_, err := service.GetNames()

		assert.Error(t, err)
	})
}

func TestGetUniqueNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		service := db.New(dbConn)
		names, err := service.GetUniqueNames()

		assert.NoError(t, err)
		assert.Contains(t, names, "Alice")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query_error", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errors.New("db error"))

		service := db.New(dbConn)
		_, err := service.GetUniqueNames()

		assert.Error(t, err)
	})

	t.Run("scan_error", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		service := db.New(dbConn)
		_, err := service.GetUniqueNames()

		assert.Error(t, err)
	})

	t.Run("rows_err_loop", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		rows := sqlmock.NewRows([]string{"name"}).RowError(0, errors.New("loop error"))
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		service := db.New(dbConn)
		_, err := service.GetUniqueNames()

		assert.Error(t, err)
	})

	t.Run("close_error", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		rows := sqlmock.NewRows([]string{"name"}).CloseError(errors.New("close error"))
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		service := db.New(dbConn)
		_, err := service.GetUniqueNames()

		assert.Error(t, err)
	})
}
