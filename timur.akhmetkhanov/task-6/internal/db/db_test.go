package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDBService_GetNames(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := New(db)

	query := "SELECT name FROM users"

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery(query).WillReturnRows(rows)

		names, err := service.GetNames()

		assert.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Bob"}, names)
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery(query).WillReturnError(errors.New("connection failed"))

		names, err := service.GetNames()

		assert.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "db query")
	})

	t.Run("rows iteration error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			RowError(1, errors.New("row failure"))

		mock.ExpectQuery(query).WillReturnRows(rows)

		names, err := service.GetNames()

		assert.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "rows error")
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := New(db)

	query := "SELECT DISTINCT name FROM users"

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Charlie").
			AddRow("Dave")

		mock.ExpectQuery(query).WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		assert.NoError(t, err)
		assert.Equal(t, []string{"Charlie", "Dave"}, names)
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery(query).WillReturnError(errors.New("db dead"))

		names, err := service.GetUniqueNames()

		assert.Error(t, err)
		assert.Nil(t, names)
	})
}
