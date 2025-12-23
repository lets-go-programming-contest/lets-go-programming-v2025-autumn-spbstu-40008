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
    rows, err := s.DB.Query("SELECT name FROM users")
    if err != nil {
        return nil, fmt.Errorf("query error: %w", err)
    }
    
    defer func() {

        _ = rows.Close()
    }()
    
    var names []string
    for rows.Next() {
        var name string
        if scanErr := rows.Scan(&name); scanErr != nil {
            return nil, fmt.Errorf("scan error: %w", scanErr)
        }
        names = append(names, name)
    }
    
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows iteration error: %w", err)
    }
    
    return names, nil
}

func (s DBService) GetUniqueNames() ([]string, error) {
    rows, err := s.DB.Query("SELECT DISTINCT name FROM users")
    if err != nil {
        return nil, fmt.Errorf("query error: %w", err)
    }
    
    defer func() {
        _ = rows.Close()
    }()
    
    var names []string
    for rows.Next() {
        var name string
        if scanErr := rows.Scan(&name); scanErr != nil {
            return nil, fmt.Errorf("scan error: %w", scanErr)
        }
        names = append(names, name)
    }
    
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows iteration error: %w", err)
    }
    
    return names, nil
}
