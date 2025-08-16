-- name: CreateUser :one
insert into users (id, created_at, updated_at, name) values (
    $1, $2, $3, $4
) on conflict (name) do nothing returning *;

-- name: GetUserByName :one
select * from users where name = $1;

-- name: GetAllUsers :many
select * from users;

-- name: DeleteAll :exec
truncate users restart identity cascade;
