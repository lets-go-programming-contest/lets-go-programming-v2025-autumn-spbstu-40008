package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type DatabaseQueryExecutor interface {
	QueryContext(ctx context.Context, sqlQuery string, params ...any) (*sql.Rows, error)
}

func extractString(scanner interface{ Scan(...any) error }, output *string) error {
	if scanErr := scanner.Scan(output); scanErr != nil {
		return fmt.Errorf("string extraction failed: %w", scanErr)
	}
	return nil
}

type nameFetcher struct {
	executor DatabaseQueryExecutor
	sqlQuery string
}

func (f *nameFetcher) Retrieve(ctx context.Context) ([]string, error) {
	rows, queryErr := f.executor.QueryContext(ctx, f.sqlQuery)
	if queryErr != nil {
		return nil, fmt.Errorf("database query execution failed: %w", queryErr)
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var nameValue string
		if extractErr := extractString(rows, &nameValue); extractErr != nil {
			return nil, extractErr
		}
		results = append(results, nameValue)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, errors.Join(fmt.Errorf("row processing failed"), rowsErr)
	}

	return results, nil
}

type UserHandler struct {
	nameRetriever NameRetriever
}

type NameRetriever interface {
	Retrieve(ctx context.Context) ([]string, error)
}

func InitializeUserHandler(db DatabaseQueryExecutor) *UserHandler {
	return &UserHandler{
		nameRetriever: &nameFetcher{
			executor: db,
			sqlQuery: "SELECT name FROM users",
		},
	}
}

func (h *UserHandler) GetAllNames(ctx context.Context) ([]string, error) {
	return h.nameRetriever.Retrieve(ctx)
}

func (h *UserHandler) GetDistinctNamesAsSet(ctx context.Context) (map[string]struct{}, error) {
	rows, queryErr := h.nameRetriever.(*nameFetcher).executor.QueryContext(ctx, "SELECT DISTINCT name FROM users")
	if queryErr != nil {
		return nil, fmt.Errorf("distinct name query failed: %w", queryErr)
	}
	defer rows.Close()

	resultSet := make(map[string]struct{})
	for rows.Next() {
		var nameValue string
		if extractErr := extractString(rows, &nameValue); extractErr != nil {
			return nil, extractErr
		}
		resultSet[nameValue] = struct{}{}
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, errors.Join(fmt.Errorf("distinct name processing failed"), rowsErr)
	}

	return resultSet, nil
}
