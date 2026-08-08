-- name: SelectPairs :many
SELECT DISTINCT
    pair_base,
    pair_quote
FROM orders;


-- name: InsertOrder :exec
INSERT INTO orders (
    id,
    user_id,
    pair_base,
    pair_quote,
    side,
    price,
    qty,
    remaining_qty,
    created_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);