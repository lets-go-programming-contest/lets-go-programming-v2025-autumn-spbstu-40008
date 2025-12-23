package db_test

import (
	"errors"
	"testing"

	"example_mock/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errQuery    = errors.New("query error")
	errRow      = errors.New("row error")
)

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		expected    []string
		expectedErr error
	}{
		{
			name: "success - return names",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Ivan").
					AddRow("Gena228")
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rows)
			},
			expected: []string{"Ivan", "Gena228"},
		},
		{
			name: "error - query error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnError(errQuery)
			},
			expectedErr: errQuery,
		},
		{
			name: "error - scan error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow(nil)
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rows)
			},
			expectedErr: errors.New("sql: Scan error"), 
		},
		{
			name: "error - rows iteration error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Ivan").
					RowError(1, errRow)
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rows)
			},
			expectedErr: errRow,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := db.New(mockDB)

			tc.setupMock(mock)

			result, err := service.GetNames()

			if tc.expectedErr != nil {
				require.Error(t, err)
				
				if tc.name == "error - scan error" {
					assert.Contains(t, err.Error(), "converting NULL to string is unsupported")
				} else {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	assert.NotNil(t, service)
	assert.Equal(t, mockDB, service.DB)
}
