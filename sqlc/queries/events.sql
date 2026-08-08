-- name: TryLockPair :one
SELECT pg_try_advisory_xact_lock(hashtext(sqlc.arg(pair_key)));

-- name: SelectOrderEventBatch :many
SELECT
    e.id,
    e.aggregate_id,
    e.aggregate_type,
    e.partition_key,
    e.event_type,
    e.payload,
    e.created_at
FROM events e
WHERE e.aggregate_type = 'order' AND
    e.status = 'pending' AND
    e.partition_key = sqlc.arg(pair_key)
ORDER BY e.id ASC
LIMIT sqlc.arg(batch_size);

-- name: ResolveEventBatch :exec
UPDATE events
SET "status" = 'resolved'
WHERE id = ANY(sqlc.slice(event_ids));

-- name: InsertEvent :exec
INSERT INTO events (
    id,
    aggregate_id,
    aggregate_type,
    partition_key,
    event_type,
    payload,
    created_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7);