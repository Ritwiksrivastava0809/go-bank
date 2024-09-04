-- name: AddEntry :exec
INSERT INTO entries(
    account_id ,
    amount
)VALUES(
    $1,$2
);

-- name: GetEntry :one
SELECT * FROM entries
WHERE account_id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
ORDER BY id 
LIMIT $1
OFFSET $2;

