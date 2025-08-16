-- name: CreateFeed :one
select * from add_feed(sqlc.arg(id)::uuid,
    sqlc.arg(created_at)::timestamp,
    sqlc.arg(updated_at)::timestamp,
    sqlc.arg(name)::text,
    sqlc.arg(url)::text,
    sqlc.arg(username)::text);

-- name: GetFeedsWithUsers :many
select feeds.name, feeds.url, users.name from users
inner join feeds on feeds.user_id = users.id;
