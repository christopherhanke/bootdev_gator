-- name: CreatePost :one
insert into posts (
    id,
    created_at,
    updated_at,
    title,
    url,
    description,
    published_at,
    feed_id
) values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
) returning *;

-- name: GetPostsForUser :many
select posts.* from posts
inner join feeds on posts.feed_id = feeds.id
inner join users on feeds.user_id = users.id
where users.name = $1
order by posts.published_at desc
limit $2;
