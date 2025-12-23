package db

import (
	"database.sql"
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

func (service DBService) GetNames() ([]string, error) {
	queryString := "SELECT name FROM users"

	userNameRows, err := service.DB.Query(queryString)
	if err != nil {
		return nil, fmt.Errorf("db query: %w", err)
	}
	defer userNameRows.Close()

	var userNameList []string

	for userNameRows.Next() {
		var userName string

		if err := userNameRows.Scan(&userName); err != nil {
			return nil, fmt.Errorf("rows scanning: %w", err)
		}

		userNameList = append(userNameList, userName)
	}

	if err := userNameRows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return userNameList, nil
}

func (service DBService) GetUniqueNames() ([]string, error) {
	queryString := "SELECT DISTINCT name FROM users"

	uniqueUserNameRows, err := service.DB.Query(queryString)
	if err != nil {
		return nil, fmt.Errorf("db query: %w", err)
	}
	defer uniqueUserNameRows.Close()

	var uniqueUserNameList []string

	for uniqueUserNameRows.Next() {
		var userName string

		if err := uniqueUserNameRows.Scan(&userName); err != nil {
			return nil, fmt.Errorf("rows scanning: %w", err)
		}

		uniqueUserNameList = append(uniqueUserNameList, userName)
	}

	if err := uniqueUserNameRows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return uniqueUserNameList, nil
}
