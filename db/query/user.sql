-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    full_name,
    email,
    role
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users 
WHERE username =$1
LIMIT 1;

-- name: UpdateUser :one
Update users 
SET 
    hashed_password =  coalesce(sqlc.narg('hashed_password'),hashed_password),
    password_changed_at =  coalesce(sqlc.narg('password_changed_at'),password_changed_at),
    full_name= coalesce(sqlc.narg('full_name'),full_name),
    email= coalesce(sqlc.narg('email'),email),
    is_email_verified = COALESCE(sqlc.narg(is_email_verified), is_email_verified)
WHERE username = @username
RETURNING *;