package db

import (
	"database/sql"
	"fmt"
)

type Database interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

type InventoryService struct {
	DB Database
}

func New(db Database) InventoryService {
	return InventoryService{DB: db}
}

func (s InventoryService) GetStockItems() ([]string, error) {
	query := "SELECT item_name FROM inventory WHERE quantity > 0"
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("inventory query: %w", err)
	}
	defer rows.Close()

	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}
		items = append(items, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return items, nil
}
