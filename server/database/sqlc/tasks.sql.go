// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: tasks.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createTask = `-- name: CreateTask :one
INSERT INTO task (
  title, description
) VALUES (
  $1, $2
)
RETURNING id, title, description
`

type CreateTaskParams struct {
	Title       string
	Description pgtype.Text
}

func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error) {
	row := q.db.QueryRow(ctx, createTask, arg.Title, arg.Description)
	var i Task
	err := row.Scan(&i.ID, &i.Title, &i.Description)
	return i, err
}

const deleteTask = `-- name: DeleteTask :exec
DELETE FROM task
WHERE id = $1
`

func (q *Queries) DeleteTask(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteTask, id)
	return err
}

const getTask = `-- name: GetTask :one
SELECT id, title, description FROM task
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTask(ctx context.Context, id pgtype.UUID) (Task, error) {
	row := q.db.QueryRow(ctx, getTask, id)
	var i Task
	err := row.Scan(&i.ID, &i.Title, &i.Description)
	return i, err
}

const listTasks = `-- name: ListTasks :many
SELECT id, title, description FROM task
ORDER BY title
`

func (q *Queries) ListTasks(ctx context.Context) ([]Task, error) {
	rows, err := q.db.Query(ctx, listTasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(&i.ID, &i.Title, &i.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTask = `-- name: UpdateTask :exec
UPDATE task
  set title = $2,
  description = $3
WHERE id = $1
`

type UpdateTaskParams struct {
	ID          pgtype.UUID
	Title       string
	Description pgtype.Text
}

func (q *Queries) UpdateTask(ctx context.Context, arg UpdateTaskParams) error {
	_, err := q.db.Exec(ctx, updateTask, arg.ID, arg.Title, arg.Description)
	return err
}
