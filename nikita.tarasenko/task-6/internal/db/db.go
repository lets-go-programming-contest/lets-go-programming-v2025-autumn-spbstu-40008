package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrQueryFailed   = errors.New("failed to query database")
	ErrScanFailed    = errors.New("failed to scan row")
	ErrIterationFail = errors.New("error during row iteration")
)

type DBQuerier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func scanString(scanner interface{ Scan(...any) error }, target *string) error {
	if err := scanner.Scan(target); err != nil {
		return fmt.Errorf("%w: %v", ErrScanFailed, err)
	}
	return nil
}

type nameCollector struct {
	querier DBQuerier
	query   string
}

func (c *nameCollector) Collect(ctx context.Context) ([]string, error) {
	rows, err := c.querier.QueryContext(ctx, c.query)
	if err != nil {
		return nil, errors.Join(ErrQueryFailed, err)
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var name string
		if err := scanString(rows, &name); err != nil {
			return nil, err
		}
		result = append(result, name)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Join(ErrIterationFail, err)
	}

	return result, nil
}

type UserService struct {
	nameLister NameCollector
}

type NameCollector interface {
	Collect(ctx context.Context) ([]string, error)
}

func NewUserService(db DBQuerier) *UserService {
	return &UserService{
		nameLister: &nameCollector{
			querier: db,
			query:   "SELECT name FROM users",
		},
	}
}

func (s *UserService) ListAllNames(ctx context.Context) ([]string, error) {
	return s.nameLister.Collect(ctx)
}

func (s *UserService) ListUniqueNamesAsSet(ctx context.Context) (map[string]struct{}, error) {
	rows, err := s.nameLister.(*nameCollector).querier.QueryContext(ctx, "SELECT DISTINCT name FROM users")
	if err != nil {
		return nil, errors.Join(ErrQueryFailed, err)
	}
	defer rows.Close()

	result := make(map[string]struct{})
	for rows.Next() {
		var name string
		if err := scanString(rows, &name); err != nil {
			return nil, err
		}
		result[name] = struct{}{}
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Join(ErrIterationFail, err)
	}

	return result, nil
}
func _() {
	_ = ErrQueryFailed
	_ = NewUserService
}
