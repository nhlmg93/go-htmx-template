-- name: GetTodo :one
SELECT * FROM todos 
WHERE id = ? LIMIT 1;

-- name: ListTodos :many
SELECT * FROM todos
ORDER BY id;

-- name: CreateTodo :one
INSERT INTO todos (
  description
) VALUES (
  ?
)
RETURNING *;

-- name: UpdateTodo :exec
UPDATE todos
set description = ?
WHERE id = ?
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = ?;
