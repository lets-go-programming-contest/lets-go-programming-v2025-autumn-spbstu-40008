package db

import (
    "database/sql"
    "errors"
    "fmt"
    "log"
)

var (
    ErrUserNotFound   = errors.New("user not found")
    ErrNoRowsAffected = errors.New("no rows affected")
    ErrQueryFailed    = errors.New("query failed")
    ErrRowsError      = errors.New("rows error")
)

type Database interface {
	Query(query string, args ...any) (*sql.Rows, error)
    QueryRow(query string, args ...any) *sql.Row
    Exec(query string, args ...any) (sql.Result, error)
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
        return nil, fmt.Errorf("%w: %w", ErrQueryFailed, err)
    }

    defer rows.Close()

    names := make([]string, 0)

    for rows.Next() {
        var name string
        if err := rows.Scan(&name); err != nil {
            log.Printf("scan error: %v", err)

            continue
        }

        names = append(names, name)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("%w: %w", ErrRowsError, err)
    }

    return names, nil
}

func (service DBService) AddUser(name string) error {
    query := "INSERT INTO users (name) VALUES ($1)"

    _, err := service.DB.Exec(query, name)
    if err != nil {
        return fmt.Errorf("add user: %w", err)
    }

    return nil
}

func (service DBService) GetUserByID(id int) (string, error) {
    var name string

    query := "SELECT name FROM users WHERE id = $1"

    err := service.DB.QueryRow(query, id).Scan(&name)
    if errors.Is(err, sql.ErrNoRows) {
        return "", fmt.Errorf("user with id %d: %w", id, ErrUserNotFound)
    }

    if err != nil {
        return "", fmt.Errorf("get user by id: %w", err)
    }

    return name, nil
}

func (service DBService) UpdateUser(id int, newName string) error {
    query := "UPDATE users SET name = $1 WHERE id = $2"

    result, err := service.DB.Exec(query, newName, id)
    if err != nil {
        return fmt.Errorf("update user: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("user with id %d: %w", id, ErrUserNotFound)
    }

    return nil
}

func (service DBService) DeleteUser(id int) error {
    query := "DELETE FROM users WHERE id = $1"

    result, err := service.DB.Exec(query, id)
    if err != nil {
        return fmt.Errorf("delete user: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("user with id %d: %w", id, ErrUserNotFound)
    }

    return nil
}
