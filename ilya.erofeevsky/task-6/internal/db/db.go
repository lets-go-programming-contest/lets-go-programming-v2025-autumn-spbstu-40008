package db

import (
	"database/sql"
)

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetNames() ([]string, error) {
	rows, err := s.db.Query("SELECT name FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []string

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		result = append(result, name)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) GetUniqueNames() ([]string, error) {
	rows, err := s.db.Query("SELECT DISTINCT name FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []string

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		result = append(result, name)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return result, nil
}

