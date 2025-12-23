package db_test

import (
	"errors"
	"testing"

	"example_mock/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBService_GetNames(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	type testCase struct {
		name          string
		mockBehavior  func()
		expectedNames []string
		expectError   bool
		errorContains string
	}

	tests := []testCase{
		{
			name: "Success",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Ivan").
					AddRow("Gena228")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedNames: []string{"Ivan", "Gena228"},
			expectError:   false,
		},
		{
			name: "Query Error",
			mockBehavior: func() {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnError(errors.New("connection failed"))
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "failed to execute query",
		},
		{
			name: "Scan Error (column mismatch)",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"name", "age"}).
					AddRow("Ivan", 25)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "failed to scan row",
		},
		{
			name: "Iteration Error (rows.Err)",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Ivan").
					RowError(0, errors.New("corrupted data"))
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "error during iteration",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			names, err := service.GetNames()

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedNames, names)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}