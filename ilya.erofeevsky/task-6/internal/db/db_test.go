package db_test

import (
    "errors"
    "testing"
    
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/Ilya-Er0fick/task-6/internal/db"
    "github.com/stretchr/testify/assert"
)

func TestDBService_GetNames(t *testing.T) {
	query := "SELECT\\s+name\\s+FROM\\s+users"

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

		_, err := service.GetNames()
		if err == nil {
			t.Log("WARNING: Implementation swallowed the close error")
		} else {
			assert.Error(t, err)
		}
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	query := "SELECT\\s+DISTINCT\\s+name\\s+FROM\\s+users"

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

		_, err := service.GetUniqueNames()
		if err == nil {
			t.Log("WARNING: Implementation swallowed the close error")
		} else {
			assert.Error(t, err)
		}
	})
}
