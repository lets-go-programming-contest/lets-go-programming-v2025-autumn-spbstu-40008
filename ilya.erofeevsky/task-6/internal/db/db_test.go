package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ilya-Er0fick/task-6/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBService_GetNames_Success(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	service := db.New(dbConn)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	names, err := service.GetNames()
	assert.NoError(t, err)
	assert.Equal(t, []string{"Alice", "Bob"}, names)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames_EmptyResult(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	service := db.New(dbConn)

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	names, err := service.GetNames()
	assert.NoError(t, err)
	assert.Empty(t, names)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames_QueryError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	service := db.New(dbConn)

	mock.ExpectQuery("SELECT name FROM users").WillReturnError(errors.New("connection lost"))

	names, err := service.GetNames()
	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "query error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames_ScanError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	service := db.New(dbConn)

	rows := sqlmock.NewRows([]string{"name"}).AddRow(123)
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	names, err := service.GetNames()
	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "scan error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames_RowsError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	service := db.New(dbConn)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errors.New("iteration error"))

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	names, err := service.GetNames()
	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "rows iteration error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames_CloseError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	service := db.New(dbConn)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		CloseError(errors.New("close failed"))

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	names, err := service.GetNames()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "close error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_Success(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	service := db.New(dbConn)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	assert.NoError(t, err)
	assert.Equal(t, []string{"Alice", "Bob"}, names)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_EmptyResult(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	service := db.New(dbConn)

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	assert.NoError(t, err)
	assert.Empty(t, names)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_QueryError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	service := db.New(dbConn)

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errors.New("connection lost"))

	names, err := service.GetUniqueNames()
	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "query error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_ScanError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	service := db.New(dbConn)

	rows := sqlmock.NewRows([]string{"name"}).AddRow(123)
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "scan error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_RowsError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	service := db.New(dbConn)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errors.New("iteration error"))

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "rows iteration error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_CloseError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	service := db.New(dbConn)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		CloseError(errors.New("close failed"))

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "close error")
	assert.NoError(t, mock.ExpectationsWereMet())
}
