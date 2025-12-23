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

func (s DBService) GetNames() ([]string, error) {
	const sqlStmt = "SELECT name FROM users"
	rows, err := s.DB.Query(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var output []string
	for rows.Next() {
		var record string
		if err := rows.Scan(&record); err != nil {
			return nil, fmt.Errorf("failed to read row data: %w", err)
		}
		output = append(output, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during iteration: %w", err)
	}

	return output, nil
}

func (s DBService) GetUniqueNames() ([]string, error) {
	const sqlStmt = "SELECT DISTINCT name FROM users"
	rows, err := s.DB.Query(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("distinct query failed: %w", err)
	}
	defer rows.Close()

	var distinct []string
	for rows.Next() {
		var item string
		if err := rows.Scan(&item); err != nil {
			return nil, fmt.Errorf("distinct row scan error: %w", err)
		}
		distinct = append(distinct, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("distinct iteration error: %w", err)
	}

	return distinct, nil
}
