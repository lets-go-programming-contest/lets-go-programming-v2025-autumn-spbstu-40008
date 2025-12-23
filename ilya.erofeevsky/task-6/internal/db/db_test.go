package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("User1").AddRow("User2")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		res, err := service.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"User1", "User2"}, res)
	})

	t.Run("query_error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := New(sqlDB)

		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errors.New("fail"))

		res, err := service.GetNames()
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "query error")
	})

	t.Run("scan_error", func(t *testing.T) {
		sqlDB, mock, _ := sqlmock.New()
		defer sqlDB.Close()
		service := New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow(123)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		res, err := service.GetNames()
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "scan error")
	})
}
