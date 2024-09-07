-- name: InsertTransaction :one
INSERT INTO transfers(
    from_account_id,
    to_account_id,
    amount
)VALUES(
    $1 , $2 , $3
)RETURNING *;

-- name: GetTransactionByID :one
SELECT * FROM transfers 
WHERE id = $1 LIMIT 1;

-- name: GetTransactionHistoryByAccountID :many
SELECT * FROM transfers
WHERE from_account_id = $1 
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: DeleteTransactionByID :exec
DELETE FROM 
transfers WHERE 
id = $1;

-- name: ListTransfers :many
SELECT * FROM transfers
WHERE 
    from_account_id = $1 OR
    to_account_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;