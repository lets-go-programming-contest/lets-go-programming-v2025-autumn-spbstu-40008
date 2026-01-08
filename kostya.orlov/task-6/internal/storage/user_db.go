package storage

import (
    "database/sql"
    "fmt"
)

type UserStore struct {
    db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
    return &UserStore{db: db}
}

func (s *UserStore) GetUserEmails() ([]string, error) {
    rows, err := s.db.Query("SELECT email FROM users")
    if err != nil {
        return nil, fmt.Errorf("db execute error: %w", err)
    }
    defer rows.Close()

    var list []string
    for rows.Next() {
        var email string
        if err := rows.Scan(&email); err != nil {
            return nil, fmt.Errorf("scan failure: %w", err)
        }
        list = append(list, email)
    }
    return list, nil
}