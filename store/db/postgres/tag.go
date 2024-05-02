package postgres

import (
	"context"
	"strings"

	"github.com/usememos/memos/store"
)

func (d *DB) UpsertTag(ctx context.Context, upsert *store.Tag) (*store.Tag, error) {
	stmt := "INSERT INTO tag (name, creator_id) VALUES ($1, $2) ON CONFLICT (name, creator_id) DO UPDATE SET name = $3"
	if _, err := d.db.ExecContext(ctx, stmt, upsert.Name, upsert.CreatorID, upsert.Name); err != nil {
		return nil, err
	}
	return upsert, nil
}

func (d *DB) ListTags(ctx context.Context, find *store.FindTag) ([]*store.Tag, error) {
	where, args := []string{"1 = 1"}, []any{}

	if find.CreatorID != 0 {
		where, args = append(where, "creator_id = "+placeholder(len(args)+1)), append(args, find.CreatorID)
	}

	query := `
		SELECT
			name,
			creator_id
		FROM tag
		WHERE ` + strings.Join(where, " AND ") + `
		ORDER BY name ASC
	`
	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*store.Tag{}
	for rows.Next() {
		tag := &store.Tag{}
		if err := rows.Scan(
			&tag.Name,
			&tag.CreatorID,
		); err != nil {
			return nil, err
		}

		list = append(list, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (d *DB) DeleteTag(ctx context.Context, delete *store.DeleteTag) error {
	where, args := []string{"name = $1", "creator_id = $2"}, []any{delete.Name, delete.CreatorID}
	stmt := `DELETE FROM tag WHERE ` + strings.Join(where, " AND ")
	result, err := d.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}
