// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: memory.sql

package db

import (
	"context"
)

const createMemory = `-- name: CreateMemory :one
INSERT INTO memories
( content, user_id )
VALUES
( $1, $2 )
RETURNING id, content, user_id, created_at, updated_at
`

type CreateMemoryParams struct {
	Content string
	UserID  int64
}

func (q *Queries) CreateMemory(ctx context.Context, arg CreateMemoryParams) (Memory, error) {
	row := q.db.QueryRow(ctx, createMemory, arg.Content, arg.UserID)
	var i Memory
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteMemory = `-- name: DeleteMemory :exec
DELETE FROM memories WHERE id = $1
`

func (q *Queries) DeleteMemory(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteMemory, id)
	return err
}

const getMemories = `-- name: GetMemories :many
SELECT id, content, user_id, created_at, updated_at FROM memories
`

func (q *Queries) GetMemories(ctx context.Context) ([]Memory, error) {
	rows, err := q.db.Query(ctx, getMemories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Memory
	for rows.Next() {
		var i Memory
		if err := rows.Scan(
			&i.ID,
			&i.Content,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
