package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mordw1n/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

var (
	errDatabase     = errors.New("database error")
	errConstraint   = errors.New("constraint violation")
	errConnection   = errors.New("connection failed")
	errForeignKey   = errors.New("foreign key violation")
	errRowsAffected = errors.New("rows affected not supported")
)

func TestGetNames(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer mockDB.Close()

	dbService := db.New(mockDB)

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
					AddRow("Mikhail").
					AddRow("Phillip").
					AddRow("Alex")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{"Mikhail", "Phillip", "Alex"},
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
					WillReturnError(errDatabase)
			},
			wantNames: nil,
			wantError: true,
			errorMsg:  "query failed",
		},
		{
			name: "scan error in one row",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Mikhail").
					AddRow(nil).
					AddRow("Alex")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{"Mikhail", "Alex"},
			wantError: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

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
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer mockDB.Close()

	dbService := db.New(mockDB)

	t.Run("successful insert", func(t *testing.T) {
		t.Parallel()

		mock.ExpectExec("INSERT INTO users \\(name\\) VALUES \\(\\$1\\)").
			WithArgs("John").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := dbService.AddUser("John")
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("insert error", func(t *testing.T) {
		t.Parallel()

		mock.ExpectExec("INSERT INTO users \\(name\\) VALUES \\(\\$1\\)").
			WithArgs("John").
			WillReturnError(errConstraint)

		err := dbService.AddUser("John")
		require.Error(t, err)
		require.Contains(t, err.Error(), "constraint violation")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("insert with empty name", func(t *testing.T) {
		t.Parallel()

		mock.ExpectExec("INSERT INTO users \\(name\\) VALUES \\(\\$1\\)").
			WithArgs("").
			WillReturnResult(sqlmock.NewResult(2, 1))

		err := dbService.AddUser("")
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUserByID(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer mockDB.Close()

	dbService := db.New(mockDB)

	t.Run("user found", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Mikhail")
		mock.ExpectQuery("SELECT name FROM users WHERE id = \\$1").
			WithArgs(1).
			WillReturnRows(rows)

		name, err := dbService.GetUserByID(1)
		require.NoError(t, err)
		require.Equal(t, "Mikhail", name)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user not found", func(t *testing.T) {
		t.Parallel()

		mock.ExpectQuery("SELECT name FROM users WHERE id = \\$1").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		name, err := dbService.GetUserByID(999)
		require.Error(t, err)
		require.Empty(t, name)
		require.Contains(t, err.Error(), "not found")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		t.Parallel()

		mock.ExpectQuery("SELECT name FROM users WHERE id = \\$1").
			WithArgs(2).
			WillReturnError(errConnection)

		name, err := dbService.GetUserByID(2)
		require.Error(t, err)
		require.Empty(t, name)
		require.Contains(t, err.Error(), "connection failed")
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer mockDB.Close()

	dbService := db.New(mockDB)

	t.Run("successful update", func(t *testing.T) {
		t.Parallel()

		mock.ExpectExec("UPDATE users SET name = \\$1 WHERE id = \\$2").
			WithArgs("NewName", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := dbService.UpdateUser(1, "NewName")
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("update non-existent user", func(t *testing.T) {
		t.Parallel()

		mock.ExpectExec("UPDATE users SET name = \\$1 WHERE id = \\$2").
			WithArgs("NewName", 999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := dbService.UpdateUser(999, "NewName")
		require.Error(t, err)
		require.Contains(t, err.Error(), "not found")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("update with database error", func(t *testing.T) {
		t.Parallel()

		mock.ExpectExec("UPDATE users SET name = \\$1 WHERE id = \\$2").
			WithArgs("NewName", 1).
			WillReturnError(errConstraint)

		err := dbService.UpdateUser(1, "NewName")
		require.Error(t, err)
		require.Contains(t, err.Error(), "constraint violation")
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer mockDB.Close()

	dbService := db.New(mockDB)

	t.Run("successful delete", func(t *testing.T) {
		t.Parallel()

		mock.ExpectExec("DELETE FROM users WHERE id = \\$1").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := dbService.DeleteUser(1)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("delete non-existent user", func(t *testing.T) {
		t.Parallel()

		mock.ExpectExec("DELETE FROM users WHERE id = \\$1").
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := dbService.DeleteUser(999)
		require.Error(t, err)
		require.Contains(t, err.Error(), "not found")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("delete with database error", func(t *testing.T) {
		t.Parallel()

		mock.ExpectExec("DELETE FROM users WHERE id = \\$1").
			WithArgs(1).
			WillReturnError(errForeignKey)

		err := dbService.DeleteUser(1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "foreign key violation")
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)

	defer mockDB.Close()

	service := db.New(mockDB)
	require.NotNil(t, service)
	require.Equal(t, mockDB, service.DB)
}

func TestCoverage(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer mockDB.Close()

	dbService := db.New(mockDB)

	result := sqlmock.NewErrorResult(errRowsAffected)
	mock.ExpectExec("UPDATE users SET name = \\$1 WHERE id = \\$2").
		WithArgs("Test", 1).
		WillReturnResult(result)

	err = dbService.UpdateUser(1, "Test")
	require.Error(t, err)
	require.Contains(t, err.Error(), "rows affected not supported")

	require.NoError(t, mock.ExpectationsWereMet())
}
