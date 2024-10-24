-- name: CreateUser :one
insert into users (id, created_at, updated_at, name)
values (
    $1,
    $2,
    $3,
    $4
)
returning *;

-- name: GetUsers :one
select * from users where name = $1 limit 1;

-- name: ResetUsers :exec
delete from users;