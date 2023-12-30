-- name: GetTask :one
SELECT * FROM task
WHERE id = $1 LIMIT 1;

-- name: ListTasks :many
SELECT * FROM task
ORDER BY title;

-- name: CreateTask :one
INSERT INTO task (
  title, description
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateTask :exec
UPDATE task
  set title = $2,
  description = $3
WHERE id = $1;

-- name: DeleteTask :exec
DELETE FROM task
WHERE id = $1;