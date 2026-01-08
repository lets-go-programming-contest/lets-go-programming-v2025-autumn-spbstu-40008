package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/task-6/internal/db"
)

func TestDBService_All(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	service := db.New(dbMock)

	t.Run("GetNames_Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("User1")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
		res, err := service.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"User1"}, res)
	})

	t.Run("GetUniqueNames_QueryError", func(t *testing.T) {
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errors.New("db error"))
		_, err := service.GetUniqueNames()
		assert.Error(t, err)
	})

	t.Run("GetNames_ScanError", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name", "age"}).AddRow("User1", 20)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
		_, err := service.GetNames()
		assert.Error(t, err)
	})
}