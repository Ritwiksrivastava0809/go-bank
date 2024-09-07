-- name: CreateAccount :one
INSERT INTO accounts (
    owner,
    balance,
    currency
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: GetAccountByOwner :one
SELECT * FROM accounts
WHERE owner = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT id, owner, balance, currency, created_at FROM accounts
ORDER BY 
  CASE 
    WHEN $3 = 'id' THEN id::text
    WHEN $3 = 'created_at' THEN created_at::text
    WHEN $3 = 'balance' THEN balance::text
    WHEN $3 = 'owner' THEN owner::text
    ELSE id::text  
  END 
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;


-- name: DeleteAccount :exec
DELETE FROM 
accounts WHERE
id = $1;