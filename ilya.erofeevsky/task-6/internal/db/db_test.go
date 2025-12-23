package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ilya-Er0fick/task-6/internal/db"
)

func TestGetNames_Success(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := db.New(dbConn)
	names, err := service.GetNames()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
	if names[0] != "Alice" || names[1] != "Bob" {
		t.Errorf("unexpected names: %v", names)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetNames_EmptyResult(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := db.New(dbConn)
	names, err := service.GetNames()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if names != nil && len(names) != 0 {
		t.Errorf("expected empty result, got %v", names)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetNames_QueryError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer dbConn.Close()

	mock.ExpectQuery("SELECT name FROM users").WillReturnError(errors.New("db error"))

	service := db.New(dbConn)
	names, err := service.GetNames()

	if err == nil {
		t.Error("expected error, got nil")
	}
	if names != nil {
		t.Errorf("expected nil names, got %v", names)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetNames_ScanError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow(123)
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := db.New(dbConn)
	names, err := service.GetNames()

	if err == nil {
		t.Error("expected error, got nil")
	}
	if names != nil {
		t.Errorf("expected nil names, got %v", names)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetNames_RowsError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errors.New("row error"))
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := db.New(dbConn)
	names, err := service.GetNames()

	if err == nil {
		t.Error("expected error, got nil")
	}
	if names != nil {
		t.Errorf("expected nil names, got %v", names)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetNames_CloseError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		CloseError(errors.New("close error"))
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := db.New(dbConn)
	_, err = service.GetNames()

	if err == nil {
		t.Error("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetUniqueNames_Success(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := db.New(dbConn)
	names, err := service.GetUniqueNames()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
	if names[0] != "Alice" || names[1] != "Bob" {
		t.Errorf("unexpected names: %v", names)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetUniqueNames_EmptyResult(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := db.New(dbConn)
	names, err := service.GetUniqueNames()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if names != nil && len(names) != 0 {
		t.Errorf("expected empty result, got %v", names)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetUniqueNames_QueryError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer dbConn.Close()

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errors.New("db error"))

	service := db.New(dbConn)
	names, err := service.GetUniqueNames()

	if err == nil {
		t.Error("expected error, got nil")
	}
	if names != nil {
		t.Errorf("expected nil names, got %v", names)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetUniqueNames_ScanError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow(123)
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := db.New(dbConn)
	names, err := service.GetUniqueNames()

	if err == nil {
		t.Error("expected error, got nil")
	}
	if names != nil {
		t.Errorf("expected nil names, got %v", names)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetUniqueNames_RowsError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errors.New("row error"))
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := db.New(dbConn)
	names, err := service.GetUniqueNames()

	if err == nil {
		t.Error("expected error, got nil")
	}
	if names != nil {
		t.Errorf("expected nil names, got %v", names)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetUniqueNames_CloseError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		CloseError(errors.New("close error"))
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := db.New(dbConn)
	_, err = service.GetUniqueNames()

	if err == nil {
		t.Error("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
