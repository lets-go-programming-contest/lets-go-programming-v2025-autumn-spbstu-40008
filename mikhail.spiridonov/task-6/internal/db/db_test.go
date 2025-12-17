package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestGetNames(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	dbService := New(mockDB)

	tests := []struct {
		name      string
		setupMock func()
		wantNames []string
		wantError bool
		errorMsg  string
	}{
		{
			name: "successful query with multiple rows",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Ivan").
					AddRow("Gena228").
					AddRow("Alice")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{"Ivan", "Gena228", "Alice"},
			wantError: false,
		},
		{
			name: "empty result",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{},
			wantError: false,
		},
		{
			name: "query error",
			setupMock: func() {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnError(errors.New("database error"))
			},
			wantNames: nil,
			wantError: true,
			errorMsg:  "query failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			names, err := dbService.GetNames()

			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAddUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	dbService := New(mockDB)

	t.Run("successful insert", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO users \\(name\\) VALUES \\(\\$1\\)").
			WithArgs("John").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := dbService.AddUser("John")
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("insert error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO users \\(name\\) VALUES \\(\\$1\\)").
			WithArgs("John").
			WillReturnError(errors.New("constraint violation"))

		err := dbService.AddUser("John")
		require.Error(t, err)
		require.Contains(t, err.Error(), "constraint violation")
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUserByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	dbService := New(mockDB)

	t.Run("user found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Ivan")
		mock.ExpectQuery("SELECT name FROM users WHERE id = \\$1").
			WithArgs(1).
			WillReturnRows(rows)

		name, err := dbService.GetUserByID(1)
		require.NoError(t, err)
		require.Equal(t, "Ivan", name)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT name FROM users WHERE id = \\$1").
			WithArgs(999).
			WillReturnError(errors.New("no rows in result set"))

		name, err := dbService.GetUserByID(999)
		require.Error(t, err)
		require.Empty(t, name)
		require.Contains(t, err.Error(), "not found")
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
