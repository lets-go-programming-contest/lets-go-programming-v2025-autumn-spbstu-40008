package usersdb

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrDatabaseQuery    = errors.New("database query execution failed")
	ErrRowScan          = errors.New("error scanning row data")
	ErrNoRowsAvailable  = errors.New("no rows available in result set")
	ErrResultProcessing = errors.New("error processing query results")
)

type QueryExecutor interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

type UserDataService struct {
	executor QueryExecutor
}

func NewUserService(executor QueryExecutor) *UserDataService {
	return &UserDataService{executor: executor}
}

func (s *UserDataService) FetchUsernames() ([]string, error) {
	const query = "SELECT name FROM users"
	
	rows, err := s.executor.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseQuery, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("%w: empty table", ErrNoRowsAvailable)
	}

	var usernames []string
	for {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrRowScan, err)
		}
		usernames = append(usernames, name)
		
		if !rows.Next() {
			break
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrResultProcessing, err)
	}

	return usernames, nil
}

func (s *UserDataService) FetchUniqueUsernames() ([]string, error) {
	const query = "SELECT DISTINCT name FROM users"
	
	rows, err := s.executor.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseQuery, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("%w: no unique records found", ErrNoRowsAvailable)
	}

	uniqueNames := make(map[string]struct{})
	var result []string
	
	for {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrRowScan, err)
		}
		
		if _, exists := uniqueNames[name]; !exists {
			uniqueNames[name] = struct{}{}
			result = append(result, name)
		}
		
		if !rows.Next() {
			break
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrResultProcessing, err)
	}

	return result, nil
}
