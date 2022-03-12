package dao

import (
	"context"
	"fmt"

	"github.com/Roukii/pock_multiplayer/internal/entity"
	"github.com/Roukii/pock_multiplayer/pkg/postgres"
)

const _defaultEntityCap = 64

// WorlDao -.
type WorlDao struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *WorlDao {
	return &WorlDao{pg}
}

// GetHistory -.
func (r *WorlDao) GetHistory(ctx context.Context) ([]entity.World, error) {
	sql, _, err := r.Builder.Select().From("world").ToSql()
	if err != nil {
		return nil, fmt.Errorf("WorlDao - GetHistory - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("WorlDao - GetHistory - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities := make([]entity.World, 0, _defaultEntityCap)

	for rows.Next() {
		e := entity.World{}

		err = rows.Scan(&e.Source, &e.Destination, &e.Original, &e.World)
		if err != nil {
			return nil, fmt.Errorf("WorlDao - GetHistory - rows.Scan: %w", err)
		}

		entities = append(entities, e)
	}

	return entities, nil
}

// Store -.
func (r *WorlDao) Store(ctx context.Context, t entity.World) error {
	sql, args, err := r.Builder.
		Insert("history").
		Columns("source, destination, original, World").
		Values(t.Source, t.Destination, t.Original, t.World).
		ToSql()
	if err != nil {
		return fmt.Errorf("WorlDao - Store - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("WorlDao - Store - r.Pool.Exec: %w", err)
	}

	return nil
}
