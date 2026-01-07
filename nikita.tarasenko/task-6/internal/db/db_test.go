package db_test

import (
	"errors"
	"testing"

	dbPkg "task-6/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errConn = errors.New("ошибка подключения")
	errRow  = errors.New("ошибка строки")
	errDB   = errors.New("ошибка базы данных")
)

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		mockFn      func(m sqlmock.Sqlmock)
		expectErr   bool
		expectedRes []string
	}{
		{
			name: "успешное получение имен",
			mockFn: func(m sqlmock.Sqlmock) {
				rs := sqlmock.NewRows([]string{"name"}).
					AddRow("Иван").
					AddRow("Петр").
					AddRow("Артем")
				m.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rs)
			},
			expectErr:   false,
			expectedRes: []string{"Иван", "Петр", "Артем"},
		},
		{
			name: "ошибка запроса",
			mockFn: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT name FROM users").
					WillReturnError(errConn)
			},
			expectErr:   true,
			expectedRes: nil,
		},
		{
			name: "ошибка сканирования",
			mockFn: func(m sqlmock.Sqlmock) {
				rs := sqlmock.NewRows([]string{"name"}).
					AddRow("Иван").
					AddRow(nil).
					AddRow("Артем")
				m.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rs)
			},
			expectErr:   true,
			expectedRes: nil,
		},
		{
			name: "ошибка итерации строк",
			mockFn: func(m sqlmock.Sqlmock) {
				rs := sqlmock.NewRows([]string{"name"}).
					AddRow("Иван").
					RowError(0, errRow)
				m.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rs)
			},
			expectErr:   true,
			expectedRes: nil,
		},
		{
			name: "пустой результат",
			mockFn: func(m sqlmock.Sqlmock) {
				rs := sqlmock.NewRows([]string{"name"})
				m.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rs)
			},
			expectErr:   false,
			expectedRes: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			conn, mock, mockErr := sqlmock.New()
			require.NoError(t, mockErr)
			defer conn.Close()

			tc.mockFn(mock)

			svc := dbPkg.New(conn)

			res, err := svc.GetNames()

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedRes, res)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name        string
		mockFn      func(m sqlmock.Sqlmock)
		expectErr   bool
		expectedRes []string
	}{
		{
			name: "успешное получение уникальных имен",
			mockFn: func(m sqlmock.Sqlmock) {
				rs := sqlmock.NewRows([]string{"name"}).
					AddRow("Иван").
					AddRow("Иван").
					AddRow("Петр").
					AddRow("Артем").
					AddRow("Петр")
				m.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(rs)
			},
			expectErr:   false,
			expectedRes: []string{"Иван", "Иван", "Петр", "Артем", "Петр"},
		},
		{
			name: "ошибка запроса уникальных имен",
			mockFn: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnError(errDB)
			},
			expectErr:   true,
			expectedRes: nil,
		},
		{
			name: "ошибка сканирования уникальных имен",
			mockFn: func(m sqlmock.Sqlmock) {
				rs := sqlmock.NewRows([]string{"name"}).
					AddRow("Иван").
					AddRow(nil).
					AddRow("Артем")
				m.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(rs)
			},
			expectErr:   true,
			expectedRes: nil,
		},
		{
			name: "ошибка строк уникальных имен",
			mockFn: func(m sqlmock.Sqlmock) {
				rs := sqlmock.NewRows([]string{"name"}).
					AddRow("Иван").
					RowError(0, errRow)
				m.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(rs)
			},
			expectErr:   true,
			expectedRes: nil,
		},
		{
			name: "пустой результат уникальных имен",
			mockFn: func(m sqlmock.Sqlmock) {
				rs := sqlmock.NewRows([]string{"name"})
				m.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(rs)
			},
			expectErr:   false,
			expectedRes: []string{},
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			t.Parallel()

			conn, mock, mockErr := sqlmock.New()
			require.NoError(t, mockErr)
			defer conn.Close()

			sc.mockFn(mock)

			svc := dbPkg.New(conn)

			res, err := svc.GetUniqueNames()

			if sc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, sc.expectedRes, res)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
