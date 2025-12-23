package db

import (
	"database/sql"
	"fmt"
)

type Database interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

type DBService struct {
	DB Database
}

func New(db Database) DBService {
	return DBService{DB: db}
}

func (s DBService) GetNames() (names []string, err error) {
	rows, queryErr := s.DB.Query("SELECT name FROM users")
	if queryErr != nil {
		return nil, fmt.Errorf("query error: %w", queryErr)
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close error: %w", closeErr)
		}
	}()

	for rows.Next() {
		var name string
		if scanErr := rows.Scan(&name); scanErr != nil {
			return nil, fmt.Errorf("scan error: %w", scanErr)
		}
		names = append(names, name)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("rows iteration error: %w", rowsErr)
	}

	return names, nil
}

func (s DBService) GetUniqueNames() (names []string, err error) {
	rows, queryErr := s.DB.Query("SELECT DISTINCT name FROM users")
	if queryErr != nil {
		return nil, fmt.Errorf("query error: %w", queryErr)
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close error: %w", closeErr)
		}
	}()

	for rows.Next() {
		var name string
		if scanErr := rows.Scan(&name); scanErr != nil {
			return nil, fmt.Errorf("scan error: %w", scanErr)
		}
		names = append(names, name)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("rows iteration error: %w", rowsErr)
	}

	return names, nil
}
