-- name: CreateUser :exec
INSERT INTO users (uuid)
VALUES (?);

-- name: GetAllUsers :many
SELECT *
FROM users;

-- name: GetUser :one
SELECT id
FROM users
WHERE uuid = ?
LIMIT 1;

-- name: GetCutLists :many
SELECT *
FROM cutLists
WHERE user_id = ?;

-- name: GetCutHistory :many
SELECT *
FROM cutHistory
WHERE lists_id = ?
ORDER BY lists_updated_at DESC;

-- name: UpdateCutListCount :exec
UPDATE cutLists
SET count = ?
WHERE id = ?;

-- name: DeleteCutHistory :exec
DELETE FROM cutHistory
WHERE lists_id = ?;

-- name: DeleteCutList :exec
DELETE FROM cutLists
WHERE id = ?;



-- name: GetEventLists :many
SELECT *
FROM eventLists
WHERE user_id = ?;

-- name: GetEventHistory :many
SELECT *
FROM eventHistory
WHERE lists_id = ?
ORDER BY lists_updated_at DESC;
