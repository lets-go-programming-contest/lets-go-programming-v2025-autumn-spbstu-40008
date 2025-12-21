package db

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrQueryExecution = errors.New("database query failed")
	ErrRowProcessing  = errors.New("row processing error")
	ErrNoRecords      = errors.New("no records found")
)

type DBExecutor interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

type DataHandler struct {
	DB DBExecutor
}

func New(db DBExecutor) DataHandler {
	return DataHandler{DB: db}
}

func (h DataHandler) GetNames() ([]string, error) {
	query := "SELECT name FROM users"
	rows, err := h.DB.Query(query)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryExecution, err)
	}

	defer rows.Close()

	var names []string

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrRowProcessing, err)
		}

		names = append(names, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRowProcessing, err)
	}

	if len(names) == 0 {
		return nil, fmt.Errorf("%w", ErrNoRecords)
	}

	return names, nil
}

func (h DataHandler) GetUniqueNames() ([]string, error) {
	query := "SELECT DISTINCT name FROM users"
	rows, err := h.DB.Query(query)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryExecution, err)
	}

	defer rows.Close()

	unique := make(map[string]struct{})

	var result []string

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrRowProcessing, err)
		}

		if _, exists := unique[name]; !exists {
			unique[name] = struct{}{}

			result = append(result, name)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRowProcessing, err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("%w", ErrNoRecords)
	}

	return result, nil
}