-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    full_name,
    email,
    password_changed_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING *;


-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CheckExistingUser :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE username = $1 OR email = $2
) AS user_exists;
