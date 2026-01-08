package storage_test  // <-- ВАЖНО: добавить _test

import (
    "errors"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/task-6/internal/storage"  // <-- Теперь это нормальный импорт
)

func TestUserStore_GetUserEmails(t *testing.T) {
    tests := []struct {
        name        string
        setupMock   func(mock sqlmock.Sqlmock)
        expected    []string
        expectError bool
        errorMsg    string
    }{
        {
            name: "success with multiple emails",
            setupMock: func(mock sqlmock.Sqlmock) {
                rows := sqlmock.NewRows([]string{"email"}).
                    AddRow("user1@example.com").
                    AddRow("user2@example.com").
                    AddRow("user3@example.com")
                mock.ExpectQuery("SELECT email FROM users").WillReturnRows(rows)
            },
            expected:    []string{"user1@example.com", "user2@example.com", "user3@example.com"},
            expectError: false,
        },
        {
            name: "success with empty result",
            setupMock: func(mock sqlmock.Sqlmock) {
                rows := sqlmock.NewRows([]string{"email"})
                mock.ExpectQuery("SELECT email FROM users").WillReturnRows(rows)
            },
            expected:    []string{},
            expectError: false,
        },
        {
            name: "query error",
            setupMock: func(mock sqlmock.Sqlmock) {
                mock.ExpectQuery("SELECT email FROM users").
                    WillReturnError(errors.New("connection failed"))
            },
            expected:    nil,
            expectError: true,
            errorMsg:    "db execute error",
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            db, mock, err := sqlmock.New()
            require.NoError(t, err)
            defer db.Close()

            store := storage.NewUserStore(db)  // <-- Используем полный путь
            tc.setupMock(mock)

            res, err := store.GetUserEmails()

            if tc.expectError {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tc.errorMsg)
                assert.Nil(t, res)
            } else {
                require.NoError(t, err)
                assert.Equal(t, tc.expected, res)
            }

            // Проверяем, что все ожидания выполнены
            assert.NoError(t, mock.ExpectationsWereMet())
        })
    }
}

func TestNewUserStore(t *testing.T) {
    db, _, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()

    store := storage.NewUserStore(db)  // <-- Используем полный путь
    assert.NotNil(t, store)
}