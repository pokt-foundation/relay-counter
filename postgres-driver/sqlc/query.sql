-- name: InsertRelayCount :exec
INSERT INTO relay_count (app_public_key, day, success, error, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (app_public_key, day) DO UPDATE
    SET success = relay_count.success + excluded.success,
        error = relay_count.error + excluded.error,
        updated_at = $6;
-- name: SelectRelayCounts :many
SELECT app_public_key, day, success, error, created_at, updated_at
FROM relay_count
WHERE day BETWEEN $1 AND $2;
