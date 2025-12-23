package db

import (
	"database/sql"
	"fmt"
)

// Database интерфейс для абстракции *sql.DB
type Database interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

type DBService struct {
	DB Database
}

func New(db Database) DBService {
	return DBService{DB: db}
}

// GetNames получает список имен из БД
func (s DBService) GetNames() ([]string, error) {
	query := "SELECT name FROM users"

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var names []string

	for rows.Next() {
		var name string
		// Ошибка сканирования (например, если типы не совпадают)
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		names = append(names, name)
	}

	// Ошибка итерации (например, разрыв соединения в процессе чтения)
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration: %w", err)
	}

	return names, nil
}