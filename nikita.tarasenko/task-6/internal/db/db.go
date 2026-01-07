package db

import (
	"database/sql"
	"fmt"
)

type Database interface {
	Query(q string, a ...any) (*sql.Rows, error)
}

type DBService struct {
	db Database
}

func New(d Database) DBService {
	return DBService{db: d}
}

func (s DBService) GetNames() ([]string, error) {
	const queryText = "SELECT name FROM users"

	rows, queryErr := s.db.Query(queryText)
	if queryErr != nil {
		return nil, fmt.Errorf("ошибка запроса: %w", queryErr)
	}
	defer rows.Close()

	var result []string

	for rows.Next() {
		var n string

		if scanErr := rows.Scan(&n); scanErr != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %w", scanErr)
		}

		result = append(result, n)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("ошибка обработки строк: %w", rowsErr)
	}

	return result, nil
}

func (s DBService) GetUniqueNames() ([]string, error) {
	const q = "SELECT DISTINCT name FROM users"

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var data []string

	for rows.Next() {
		var val string

		if e := rows.Scan(&val); e != nil {
			return nil, fmt.Errorf("ошибка сканирования: %w", e)
		}

		data = append(data, val)
	}

	if e := rows.Err(); e != nil {
		return nil, fmt.Errorf("ошибка итерации: %w", e)
	}

	return data, nil
}
