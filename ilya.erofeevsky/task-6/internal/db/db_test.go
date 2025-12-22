package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ilya-Er0fick/task-6/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInventoryService_GetStockItems(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)
		rows := sqlmock.NewRows([]string{"item_name"}).
			AddRow("Processor").
			AddRow("RAM")

		mock.ExpectQuery("SELECT item_name FROM inventory WHERE quantity > 0").
			WillReturnRows(rows)

		res, err := service.GetStockItems()
		assert.NoError(t, err)
		assert.Equal(t, []string{"Processor", "RAM"}, res)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db_error", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)
		mock.ExpectQuery("SELECT item_name FROM inventory").
			WillReturnError(errors.New("sql error"))

		res, err := service.GetStockItems()
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
