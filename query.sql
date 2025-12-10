-- name: CreateUser :exec
INSERT INTO users (uuid)
VALUES (?);

-- name: CreateCutList :exec
INSERT INTO cutLists (this_is_cut, user_id, name, color, count, `limit`, late_time, late_count)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name : UpdateCutList :exec
UPDATE cutLists
SET this_is_cut = ?, name = ?, color = ?,count = ? `limit` = ?, late_time = ?, late_count = ?
WHERE id = ?;

-- name: CreateCutHistory :exec
INSERT INTO cutHistory (this_is_cut, late_time, lists_id, lists_updated_at)
VALUES (?, ?, ?, ?);

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
