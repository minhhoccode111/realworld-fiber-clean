package persistent

import (
	"context"
	"fmt"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/postgres"
)

// TranslationCloneRepo -.
type TranslationCloneRepo struct {
	*postgres.Postgres
}

// NewClone -.
func NewClone(pg *postgres.Postgres) *TranslationCloneRepo {
	return &TranslationCloneRepo{pg}
}

// GetHistoryClone -.
func (r *TranslationCloneRepo) GetHistoryClone(
	ctx context.Context,
	limit, offset uint64,
) (translations []entity.TranslationClone, total uint64, err error) {
	sql, _, err := r.Builder.
		Select("source, destination, original, translation, count(*) over()").
		From("history").
		Limit(limit).
		Offset(offset).
		ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("TranslationCloneRepo - GetHistoryClone - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, 0, fmt.Errorf("TranslationCloneRepo - GetHistoryClone - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities := make([]entity.TranslationClone, 0)
	for rows.Next() {
		e := entity.TranslationClone{}
		err := rows.Scan(&e.Source, &e.Destination, &e.Original, &e.Translation, &total)
		if err != nil {
			return nil, 0, fmt.Errorf("TranslationCloneRepo - GetHistoryClone - rows.Scan: %w", err)
		}

		entities = append(entities, e)
	}

	return entities, total, nil
}

// StoreTranslation -.
func (r *TranslationCloneRepo) StoreTranslation(
	ctx context.Context,
	t entity.TranslationClone,
) error {
	sql, args, err := r.Builder.
		Insert("history").
		Columns("source, destination, original, translation").
		Values(t.Source, t.Destination, t.Original, t.Translation).
		ToSql()
	if err != nil {
		return fmt.Errorf("TranslationCloneRepo - StoreTranslation - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("TranslationCloneRepo - StoreTranslation - r.Pool.Exec: %w", err)
	}
	return nil
}
