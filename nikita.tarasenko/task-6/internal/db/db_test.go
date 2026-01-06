package db_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"task-6/internal/db"
)

func TestUserService_ListAllNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		mock  func(sqlmock.Sqlmock)
		check func(*testing.T, []string, error)
	}{
		{
			name: "returns correct list",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`^SELECT name FROM users$`).
					WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob"))
			},
			check: func(t *testing.T, result []string, err error) {
				require.NoError(t, err)
				require.ElementsMatch(t, []string{"Alice", "Bob"}, result)
			},
		},
		{
			name: "propagates query error",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`^SELECT name FROM users$`).
					WillReturnError(errors.New("io timeout"))
			},
			check: func(t *testing.T, _ []string, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "database query failed")
			},
		},
		{
			name: "returns empty on empty result",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`^SELECT name FROM users$`).
					WillReturnRows(sqlmock.NewRows([]string{"name"}))
			},
			check: func(t *testing.T, result []string, err error) {
				require.NoError(t, err)
				require.Empty(t, result)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dbMock, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbMock.Close()

			tc.mock(mock)

			service := db.NewUserService(dbMock)
			names, err := service.ListAllNames(context.Background())
			tc.check(t, names, err)

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserService_ListUniqueNamesAsSet(t *testing.T) {
	t.Parallel()

	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob").AddRow("Alice"))

	service := db.NewUserService(dbMock)
	result, err := service.ListUniqueNamesAsSet(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Contains(t, result, "Alice")
	require.Contains(t, result, "Bob")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_ListUniqueNamesAsSet_Error(t *testing.T) {
	t.Parallel()

	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
		WillReturnError(errors.New("connection lost"))

	service := db.NewUserService(dbMock)
	_, err = service.ListUniqueNamesAsSet(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), "database query for unique names failed")
	require.NoError(t, mock.ExpectationsWereMet())
}
