package db_test

import (
	"database/sql"
	"errors"
	"testing"

	dbpkg "github.com/narumov-diyar/task-6/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errNetworkHiccup = errors.New("network hiccup")

//nolint:ireturn
func newMockService() (*sql.DB, sqlmock.Sqlmock, dbpkg.DBService) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	service := dbpkg.New(db)

	return db, mock, service
}

func TestGetNamesSuccess(t *testing.T) {
	t.Parallel()

	db, mock, service := newMockService()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Emma").
		AddRow("Liam").
		AddRow("Olivia")

	mock.ExpectQuery(`^SELECT name FROM users$`).
		WillReturnRows(rows)

	names, err := service.GetNames()
	require.NoError(t, err)
	assert.Equal(t, []string{"Emma", "Liam", "Olivia"}, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNamesEmptyResult(t *testing.T) {
	t.Parallel()

	db, mock, service := newMockService()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"})

	mock.ExpectQuery(`^SELECT name FROM users$`).
		WillReturnRows(rows)

	names, err := service.GetNames()
	require.NoError(t, err)
	assert.Empty(t, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNamesQueryError(t *testing.T) {
	t.Parallel()

	db, mock, service := newMockService()
	defer db.Close()

	mock.ExpectQuery(`^SELECT name FROM users$`).
		WillReturnError(sql.ErrConnDone)

	_, err := service.GetNames()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "db query")
	require.ErrorIs(t, err, sql.ErrConnDone)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNamesScanError(t *testing.T) {
	t.Parallel()

	db, mock, service := newMockService()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

	mock.ExpectQuery(`^SELECT name FROM users$`).
		WillReturnRows(rows)

	_, err := service.GetNames()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rows scanning")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNamesRowsNextError(t *testing.T) {
	t.Parallel()

	db, mock, service := newMockService()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Noah").
		AddRow("Ava").
		RowError(1, errNetworkHiccup)

	mock.ExpectQuery(`^SELECT name FROM users$`).
		WillReturnRows(rows)

	_, err := service.GetNames()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rows error")
	assert.Contains(t, err.Error(), "network hiccup")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNamesSuccess(t *testing.T) {
	t.Parallel()

	db, mock, service := newMockService()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Ethan").
		AddRow("Sophia")

	mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
		WillReturnRows(rows)

	values, err := service.GetUniqueNames()
	require.NoError(t, err)
	assert.Equal(t, []string{"Ethan", "Sophia"}, values)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNamesEmptyResult(t *testing.T) {
	t.Parallel()

	db, mock, service := newMockService()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"})

	mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
		WillReturnRows(rows)

	values, err := service.GetUniqueNames()
	require.NoError(t, err)
	assert.Empty(t, values)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNamesQueryError(t *testing.T) {
	t.Parallel()

	db, mock, service := newMockService()
	defer db.Close()

	mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
		WillReturnError(sql.ErrConnDone)

	_, err := service.GetUniqueNames()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "db query")
	require.ErrorIs(t, err, sql.ErrConnDone)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNamesScanError(t *testing.T) {
	t.Parallel()

	db, mock, service := newMockService()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

	mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
		WillReturnRows(rows)

	_, err := service.GetUniqueNames()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rows scanning")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNamesRowsErrAfterNext(t *testing.T) {
	t.Parallel()

	db, mock, service := newMockService()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Mason").
		AddRow("Isabella").
		RowError(1, errNetworkHiccup)

	mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
		WillReturnRows(rows)

	_, err := service.GetUniqueNames()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rows error")
	assert.Contains(t, err.Error(), "network hiccup")

	assert.NoError(t, mock.ExpectationsWereMet())
}
