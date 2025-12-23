package db

import (
	"database/sql"
	"fmt"
)

// Database описывает интерфейс для выполнения запросов.
type Database interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

type DBService struct {
	DB Database
}

func New(db Database) DBService {
	return DBService{DB: db}
}

// GetNames возвращает список всех имен.
func (s DBService) GetNames() ([]string, error) {
	query := "SELECT name FROM users"

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("execution error: %w", err)
	}
	defer rows.Close()

	var names []string

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		names = append(names, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iteration error: %w", err)
	}

	return names, nil
}

// GetUniqueNames возвращает список уникальных имен.
func (s DBService) GetUniqueNames() ([]string, error) {
	query := "SELECT DISTINCT name FROM users"

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("execution error: %w", err)
	}
	defer rows.Close()

	var result []string

	for rows.Next() {
		var val string
		if err := rows.Scan(&val); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		result = append(result, val)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iteration error: %w", err)
	}

	return result, nil
}
