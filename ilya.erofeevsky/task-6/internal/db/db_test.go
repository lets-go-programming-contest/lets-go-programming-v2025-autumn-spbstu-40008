package db_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ilya-Er0fick/task-6/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestDBService_GetNames(t *testing.T) {
	query := regexp.QuoteMeta("SELECT name FROM users")

	t.Run("success", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("User1")
		mock.ExpectQuery(query).WillReturnRows(rows)

		res, err := service.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"User1"}, res)
	})

	t.Run("success empty result", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"})
		mock.ExpectQuery(query).WillReturnRows(rows)

		res, err := service.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{}, res)
	})

	t.Run("query_error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		mock.ExpectQuery(query).WillReturnError(errors.New("query failed"))

		res, err := service.GetNames()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "query error")
		assert.Nil(t, res)
	})

	t.Run("scan_error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery(query).WillReturnRows(rows)

		res, err := service.GetNames()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "scan error")
		assert.Nil(t, res)
	})

	t.Run("rows_err", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("User1").
			RowError(0, errors.New("row error"))
		mock.ExpectQuery(query).WillReturnRows(rows)

		res, err := service.GetNames()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rows iteration error")
		assert.Nil(t, res)
	})

	t.Run("close_error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("User1").
			CloseError(errors.New("close error"))
		mock.ExpectQuery(query).WillReturnRows(rows)

		res, err := service.GetNames()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "close error")
		assert.Equal(t, []string{"User1"}, res)
	})

	t.Run("scan_error_with_close_error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow(nil).
			CloseError(errors.New("close error"))
		mock.ExpectQuery(query).WillReturnRows(rows)

		res, err := service.GetNames()
		assert.Error(t, err)
		// Должна быть обе ошибки: scan error и close error
		assert.Contains(t, err.Error(), "scan error")
		assert.Contains(t, err.Error(), "close error")
		assert.Nil(t, res)
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	query := regexp.QuoteMeta("SELECT DISTINCT name FROM users")

	t.Run("success", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("User1")
		mock.ExpectQuery(query).WillReturnRows(rows)

		res, err := service.GetUniqueNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"User1"}, res)
	})

	t.Run("success empty result", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"})
		mock.ExpectQuery(query).WillReturnRows(rows)

		res, err := service.GetUniqueNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{}, res)
	})

	t.Run("query_error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		mock.ExpectQuery(query).WillReturnError(errors.New("query failed"))

		res, err := service.GetUniqueNames()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "query error")
		assert.Nil(t, res)
	})

	t.Run("scan_error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery(query).WillReturnRows(rows)

		res, err := service.GetUniqueNames()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "scan error")
		assert.Nil(t, res)
	})

	t.Run("rows_err", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("User1").
			RowError(0, errors.New("row error"))
		mock.ExpectQuery(query).WillReturnRows(rows)

		res, err := service.GetUniqueNames()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rows iteration error")
		assert.Nil(t, res)
	})

	t.Run("close_error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("User1").
			CloseError(errors.New("close error"))
		mock.ExpectQuery(query).WillReturnRows(rows)

		res, err := service.GetUniqueNames()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "close error")
		assert.Equal(t, []string{"User1"}, res)
	})

	t.Run("rows_err_with_close_error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("User1").
			RowError(0, errors.New("row error")).
			CloseError(errors.New("close error"))
		mock.ExpectQuery(query).WillReturnRows(rows)

		res, err := service.GetUniqueNames()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rows iteration error")
		assert.Contains(t, err.Error(), "close error")
		assert.Nil(t, res)
	})
}
