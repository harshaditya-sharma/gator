-- name: CreatePostLike :one
INSERT INTO post_likes(user_id, post_id, created_at)
VALUES($1, $2, $3)
RETURNING *;

-- name: DeletePostLike :exec
DELETE FROM post_likes WHERE user_id = $1 AND post_id = $2;

-- name: GetLikedPostsForUser :many
SELECT posts.*, feeds.name as feed_name FROM posts
JOIN post_likes ON posts.id = post_likes.post_id
JOIN feeds ON posts.feed_id = feeds.id
WHERE post_likes.user_id = $1
ORDER BY post_likes.created_at DESC
LIMIT $2 OFFSET $3;
