-- name: CreateFeed :one
insert into feeds(id, created_at, updated_at, name, url, user_id)
values(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
returning *;

-- name: GetFeeds :many
select feeds.name, feeds.url, users.name from feeds inner join users on feeds.user_id = users.id;
