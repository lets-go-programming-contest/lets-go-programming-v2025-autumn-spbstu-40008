package db

import (
	"database/sql"
	"fmt"
	"log"
)

type Database interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type DBService struct {
	DB Database
}

func New(db Database) DBService {
	return DBService{DB: db}
}

func (service DBService) GetNames() ([]string, error) {
	query := "SELECT name FROM users"

	rows, err := service.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Printf("scan error: %v", err)
			continue
		}
		names = append(names, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return names, nil
}

func (service DBService) AddUser(name string) error {
	query := "INSERT INTO users (name) VALUES ($1)"
	_, err := service.DB.Exec(query, name)
	return err
}

func (service DBService) GetUserByID(id int) (string, error) {
	var name string
	query := "SELECT name FROM users WHERE id = $1"
	err := service.DB.QueryRow(query, id).Scan(&name)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("user with id %d not found", id)
	}
	return name, err
}
