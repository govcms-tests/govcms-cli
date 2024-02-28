-- name: CreateInstallation :one
INSERT INTO installations (
    name, path, type
) VALUES (
    ?, ?, ?
)
RETURNING *;

-- name: DeleteInstallation :exec
DELETE FROM installations
where name = ?;

-- name: DeletePath :exec
DELETE FROM installations
where path = ?;

-- name: ListInstallations :many
SELECT * FROM installations
ORDER BY name;

-- name: ListPaths :many
SELECT path FROM installations;

-- name: GetInstallation :one
SELECT * FROM installations
WHERE name = ? LIMIT 1;

-- name: GetPath :one
SELECT path FROM  installations
WHERE name = ? LIMIT 1;
