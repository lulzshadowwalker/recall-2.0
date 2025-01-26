-- name: GetMemories :many
SELECT * FROM memories;

-- name: CreateMemory :one
INSERT INTO memories
( content, user_id )
VALUES
( $1, $2 )
RETURNING *;

-- name: DeleteMemory :exec
DELETE FROM memories WHERE id = $1;
