-- +goose Up
CREATE TABLE feed_follows (
  id UUID PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  user_id UUID NOT NULL,
  feed_id UUID NOT NULL,
  CONSTRAINT fk_feed_follows_user
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT fk_feed_follows_feed
    FOREIGN KEY (feed_id) REFERENCES feeds (id) ON DELETE CASCADE,
  CONSTRAINT unique_user_feed_pairs UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE IF EXISTS feed_follows;
