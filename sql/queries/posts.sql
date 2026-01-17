-- name: CreatePost :one
INSERT INTO posts(id, created_at, updated_at, title, url, description, published_at, feed_id) VALUES(
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8
) RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.*, feeds.name as feed_name FROM
posts JOIN feed_follows ON posts.feed_id = feed_follows.feed_id 
JOIN feeds ON posts.feed_id=feeds.id WHERE feed_follows.user_id= $1
ORDER BY published_at DESC LIMIT $2 OFFSET $3;

-- name: GetPostsForUserAsc :many
SELECT posts.*, feeds.name as feed_name FROM
posts JOIN feed_follows ON posts.feed_id = feed_follows.feed_id 
JOIN feeds ON posts.feed_id=feeds.id WHERE feed_follows.user_id= $1
ORDER BY published_at ASC LIMIT $2 OFFSET $3;

-- name: GetPostsForUserByFeed :many
SELECT posts.*, feeds.name as feed_name FROM
posts JOIN feed_follows ON posts.feed_id = feed_follows.feed_id 
JOIN feeds ON posts.feed_id=feeds.id WHERE feed_follows.user_id= $1 AND feeds.name = $3
ORDER BY published_at DESC LIMIT $2 OFFSET $4;

-- name: GetPostsForUserByFeedAsc :many
SELECT posts.*, feeds.name as feed_name FROM
posts JOIN feed_follows ON posts.feed_id = feed_follows.feed_id 
JOIN feeds ON posts.feed_id=feeds.id WHERE feed_follows.user_id= $1 AND feeds.name = $3
ORDER BY published_at ASC LIMIT $2 OFFSET $4;

-- name: GetPostsForUserMatching :many
SELECT posts.id, posts.created_at, posts.updated_at, posts.title, posts.url, posts.description, posts.published_at, posts.feed_id, feeds.name as feed_name FROM
posts JOIN feed_follows ON posts.feed_id = feed_follows.feed_id 
JOIN feeds ON posts.feed_id=feeds.id WHERE feed_follows.user_id= $1
AND (posts.title ILIKE '%' || sqlc.arg(search_term) || '%' OR posts.description ILIKE '%' || sqlc.arg(search_term) || '%')
ORDER BY published_at DESC LIMIT $2;

-- name: GetPostByURL :one
SELECT * FROM posts WHERE url = $1;