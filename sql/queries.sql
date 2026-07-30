-- name: InsertOrder :exec
INSERT INTO
    orders (
        id,
        user_id,
        pair_base,
        pair_quote,
        side,
        price,
        amount,
        remaining,
        created_at
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: InsertEvent :exec
INSERT INTO
    outbox (
        id,
        aggregate_id,
        aggregate_type,
        event_type,
        payload,
        created_at,
        is_commited
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7);