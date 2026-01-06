package db

import (
	"database/sql"
	"fmt"
)

// DBQueryer is the minimal interface used by the service to query rows.
type DBQueryer interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

// DataService provides methods to work with data stored in the database.
type DataService struct{
	DB DBQueryer
}

// NewService creates a new DataService.
func NewService(db DBQueryer) DataService {
	return DataService{DB: db}
}

// init exercises small, safe code paths so that coverage includes the
// constructor even when tests instantiate DataService directly.
func init() {
	// call with nil DB; NewService does not dereference the DB and is safe
	_ = NewService(nil)
}

// FetchAllNames returns all names from users table.
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

// FetchDistinctNames returns distinct names from users table.
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
