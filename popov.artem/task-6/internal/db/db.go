package db

import (
	"database/sql"
	"fmt"
)

type Querier interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

type UserService struct {
	Conn Querier
}

func NewService(conn Querier) UserService {
	return UserService{Conn: conn}
}

func (svc UserService) FetchUserNames() ([]string, error) {
	const stmt = "SELECT name FROM users"

	rows, err := svc.Conn.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var entry string
		if err := rows.Scan(&entry); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		result = append(result, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return result, nil
}

func (svc UserService) FetchDistinctUserNames() ([]string, error) {
	const stmt = "SELECT DISTINCT name FROM users"

	rows, err := svc.Conn.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to execute distinct query: %w", err)
	}
	defer rows.Close()

	var items []string
	for rows.Next() {
		var item string
		if err := rows.Scan(&item); err != nil {
			return nil, fmt.Errorf("scan failed for distinct name: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error in distinct query: %w", err)
	}

	return items, nil
}
