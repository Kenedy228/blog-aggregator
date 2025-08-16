-- name: CreateFeedFollow :one
select (f).id::uuid as id,
(f).created_at::timestamp as created_at,
(f).updated_at::timestamp as updated_at,
(f).user_name::text as user_name,
(f).feed_name::text as feed_name from follow_feed(
    sqlc.arg(id)::uuid, 
    sqlc.arg(created_at)::timestamp, sqlc.arg(updated_at)::timestamp,
    sqlc.arg(user_name)::text, sqlc.arg(url)::text) as f;

-- name: GetFeedFollowsForUser :many
select (f).user_name::text as user_name, 
(f).feed_name::text as feed_name, (f).feed_url::text as feed_url
from feed_follows_for_users(sqlc.arg(user_name)::text) as f;

-- name: DeleteFeedFollowsForUserByUrl :exec
select unfollow_feed(sqlc.arg(url)::text, sqlc.arg(user_name)::text);
