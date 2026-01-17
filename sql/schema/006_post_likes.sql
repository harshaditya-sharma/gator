-- +goose Up
CREATE TABLE post_likes (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY(user_id, post_id)
);

-- +goose Down
DROP TABLE post_likes;
