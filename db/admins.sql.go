// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: admins.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getAdmin = `-- name: GetAdmin :one
SELECT id, username FROM admins
WHERE id = $1
`

type GetAdminRow struct {
	ID       int32       `json:"id"`
	Username pgtype.Text `json:"username"`
}

func (q *Queries) GetAdmin(ctx context.Context, id int32) (GetAdminRow, error) {
	row := q.db.QueryRow(ctx, getAdmin, id)
	var i GetAdminRow
	err := row.Scan(&i.ID, &i.Username)
	return i, err
}

const listAdmins = `-- name: ListAdmins :many
SELECT id, username FROM admins
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListAdminsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListAdminsRow struct {
	ID       int32       `json:"id"`
	Username pgtype.Text `json:"username"`
}

func (q *Queries) ListAdmins(ctx context.Context, arg ListAdminsParams) ([]ListAdminsRow, error) {
	rows, err := q.db.Query(ctx, listAdmins, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListAdminsRow
	for rows.Next() {
		var i ListAdminsRow
		if err := rows.Scan(&i.ID, &i.Username); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
