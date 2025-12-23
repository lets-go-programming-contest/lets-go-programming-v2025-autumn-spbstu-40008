package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/task-6/internal/db"
)

func TestDBService_GetNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("User1")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		res, err := service.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"User1"}, res)
	})

	t.Run("query_error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errors.New("query failed"))

		res, err := service.GetNames()
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("scan_error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		res, err := service.GetNames()
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("rows_err", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("User1").RowError(0, errors.New("row error"))
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		res, err := service.GetNames()
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("close_error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("User1").CloseError(errors.New("close error"))
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		res, err := service.GetNames()
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("User1")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		res, err := service.GetUniqueNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"User1"}, res)
	})

	t.Run("query_error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errors.New("query failed"))

		res, err := service.GetUniqueNames()
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("scan_error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		res, err := service.GetUniqueNames()
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("rows_err", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("User1").RowError(0, errors.New("row error"))
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		res, err := service.GetUniqueNames()
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("close_error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("User1").CloseError(errors.New("close error"))
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		res, err := service.GetUniqueNames()
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
