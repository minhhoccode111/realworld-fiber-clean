package persistent

import (
	"context"
	"fmt"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/postgres"
)

// TranslationRepo -.
type TagRepo struct {
	*postgres.Postgres
}

// NewTagRepo -.
func NewTagRepo(pg *postgres.Postgres) *TagRepo {
	return &TagRepo{pg}
}

func (r *TagRepo) GetList(ctx context.Context, limit, offset uint64,
) ([]entity.TagName, uint64, error) {
	sql, _, err := r.Builder.
		Select("distinct name", "count(*) over()").
		From("tags").
		Limit(limit).
		Offset(offset).
		ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("TagRepo - GetTags - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, 0, fmt.Errorf("TagRepo - GetTags - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	tags := make([]entity.TagName, 0)

	var total uint64

	for rows.Next() {
		var name entity.TagName

		err = rows.Scan(&name, &total)
		if err != nil {
			return nil, 0, fmt.Errorf("TagRepo - GetTags - rows.Scan: %w", err)
		}

		tags = append(tags, name)
	}

	return tags, total, nil
}
