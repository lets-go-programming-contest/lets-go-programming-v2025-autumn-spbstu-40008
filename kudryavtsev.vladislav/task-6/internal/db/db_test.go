package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errDbStub = errors.New("stub database error")

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedData  []string
		shouldFail    bool
		errSubstring  string
	}{
		{
			testName: "Success path",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Ivan").
					AddRow("Maria")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedData: []string{"Ivan", "Maria"},
			shouldFail:   false,
		},
		{
			testName: "Query execution error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnError(errDbStub)
			},
			expectedData: nil,
			shouldFail:   true,
			errSubstring: "failed to execute query",
		},
		{
			testName: "Row scan error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Возвращаем NULL, который нельзя засканить в string, чтобы вызвать ошибку Scan
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedData: nil,
			shouldFail:   true,
			errSubstring: "failed to scan row",
		},
		{
			testName: "Iteration error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Ivan").
					RowError(0, errDbStub) // Имитация ошибки rows.Err()
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedData: nil,
			shouldFail:   true,
			errSubstring: "error during iteration",
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// Создаем мок
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := New(mockDB)

			// Настраиваем поведение мока
			tc.mockSetup(mock)

			// Выполняем тестируемый метод
			result, err := service.GetNames()

			// Проверки
			if tc.shouldFail {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errSubstring)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedData, result)
			}

			// Убеждаемся, что все ожидания мока выполнены
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}