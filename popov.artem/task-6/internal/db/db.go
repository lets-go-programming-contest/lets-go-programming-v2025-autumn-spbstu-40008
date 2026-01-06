package db

import (
	"database/sql"
	"fmt"
)

type DBQueryer interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

type DataService struct {
	DB DBQueryer
}

func NewService(db DBQueryer) DataService {
	return DataService{DB: db}
}

func (svc DataService) FetchAllNames() ([]string, error) {
	const query = "SELECT name FROM users"

	rows, err := svc.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("db query failed: %w", err)
	}
	defer rows.Close()

	var names []string

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scan row error: %w", err)
		}
		names = append(names, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return names, nil
}

func (svc DataService) FetchDistinctNames() ([]string, error) {
	const query = "SELECT DISTINCT name FROM users"

	rows, err := svc.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("db query failed: %w", err)
	}
	defer rows.Close()

	var uniqueNames []string

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scan row error: %w", err)
		}
		uniqueNames = append(uniqueNames, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return uniqueNames, nil
}
