package db_test

import (
	"errors"
	"regexp"
	"testing"
	

	"github.com/Ilya-Er0fick/task-6/internal/db" 
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
)

func TestDBService_GetNames(t *testing.T) {
	query := "SELECT name FROM users"

	t.Run("success", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("User1").AddRow("User2")
		mock.ExpectQuery(query).WillReturnRows(rows)

		res, err := service.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"User1", "User2"}, res)
	})

	t.Run("query_error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		mock.ExpectQuery(query).WillReturnError(errors.New("query failed"))

		res, err := service.GetNames()
		assert.Error(t, err)
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
		assert.Nil(t, res)
	})

	t.Run("rows_iteration_error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("User1").
			RowError(0, errors.New("iteration error"))
		mock.ExpectQuery(query).WillReturnRows(rows)

		res, err := service.GetNames()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rows iteration error")
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	query := "SELECT DISTINCT name FROM users"

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

	t.Run("query_error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := db.New(sqlDB)

		mock.ExpectQuery(query).WillReturnError(errors.New("query failed"))

		res, err := service.GetUniqueNames()
		assert.Error(t, err)
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
	})
}
