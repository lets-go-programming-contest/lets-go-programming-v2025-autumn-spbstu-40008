package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDBService_GetNames тестирует первую функцию GetNames
func TestDBService_GetNames(t *testing.T) {
	// Создаем мок БД
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Инициализируем сервис
	service := New(db)

	tests := []struct {
		name          string
		mockBehavior  func()
		expectedNames []string
		expectError   bool
		errorContains string
	}{
		{
			name: "Success",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedNames: []string{"Alice", "Bob"},
			expectError:   false,
		},
		{
			name: "Query Error",
			mockBehavior: func() {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnError(errors.New("connection refused"))
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "db query",
		},
		{
			name: "Rows Scan Error",
			mockBehavior: func() {
				// Создаем ситуацию ошибки сканирования: возвращаем число вместо строки или NULL, если драйвер строгий,
				// но надежнее вернуть несовместимые типы колонок, если sqlmock позволяет.
				// Проще всего: вернуть NULL для NotNull поля или просто ошибку на уровне Scan
				// В sqlmock для Scan error проще всего передать меньше колонок или несовместимые типы.
				// Однако, sqlmock хранит значения как Driver.Value.
				// Эмулируем ошибку добавлением NULL для string scan (если библиотека требует string).
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "rows scanning",
		},
		{
			name: "Rows Iteration Error",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					RowError(0, errors.New("row corrupted"))
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "rows error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			names, err := service.GetNames()

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedNames, names)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// TestDBService_GetUniqueNames тестирует вторую функцию GetUniqueNames
func TestDBService_GetUniqueNames(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	service := New(db)

	tests := []struct {
		name          string
		mockBehavior  func()
		expectedNames []string
		expectError   bool
		errorContains string
	}{
		{
			name: "Success",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
				// Важно: sqlmock использует регулярные выражения. DISTINCT нужно экранировать или писать точно.
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedNames: []string{"Alice"},
			expectError:   false,
		},
		{
			name: "Query Error",
			mockBehavior: func() {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnError(errors.New("fail"))
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "db query",
		},
		{
			name: "Rows Scan Error",
			mockBehavior: func() {
				// Передаем nil, чтобы вызвать ошибку сканирования в string
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "rows scanning",
		},
		{
			name: "Rows Iteration Error",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					RowError(0, errors.New("iteration error"))
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "rows error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()
			names, err := service.GetUniqueNames()

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedNames, names)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}