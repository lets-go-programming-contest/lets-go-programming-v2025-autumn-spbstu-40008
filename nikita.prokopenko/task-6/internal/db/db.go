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

type DataHandler struct {
	DB *sql.DB
}

func CreateHandler(db *sql.DB) *DataHandler {
	return &DataHandler{DB: db}
}

func (h *DataHandler) RetrieveNames() ([]string, error) {
	const query = "SELECT name FROM users"
	rows, err := h.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryExecution, err)
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, fmt.Errorf("%w: empty result set", ErrNoRecords)
	}
	var names []string
	for {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrRowProcessing, err)
		}
		names = append(names, name)
		if !rows.Next() {
			break
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRowProcessing, err)
	}
	return names, nil
}

func (h *DataHandler) RetrieveUniqueNames() ([]string, error) {
	const query = "SELECT DISTINCT name FROM users"
	rows, err := h.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryExecution, err)
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, fmt.Errorf("%w: no distinct records", ErrNoRecords)
	}
	unique := make(map[string]struct{})
	var result []string
	for {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrRowProcessing, err)
		}
		if _, exists := unique[name]; !exists {
			unique[name] = struct{}{}
			result = append(result, name)
		}
		if !rows.Next() {
			break
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRowProcessing, err)
	}
	return result, nil
}
