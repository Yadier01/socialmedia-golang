CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(50) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  following_count INT DEFAULT 0,
  follower_count INT DEFAULT 0,
  email VARCHAR(254) UNIQUE NOT NULL,
  created_at TIMESTAMP(0) 
  WITH
    TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS profiles (
  user_id BIGINT PRIMARY KEY,
  bio TEXT,
  avatar_url TEXT,
  join_at TIMESTAMP(0)
  WITH
    TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS posts (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  body TEXT NOT NULL,
  likes BIGINT NOT NULL DEFAULT 0,
  comments BIGINT NOT NULL DEFAULT 0,
  parent_post_id BIGINT, -- if the post is a comment 
  created_at TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  FOREIGN KEY (parent_post_id) REFERENCES posts (id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS followers (
  id BIGSERIAL PRIMARY KEY,
  follower_id BIGINT NOT NULL,
  following_id BIGINT NOT NULL,
  created_at TIMESTAMP(0)
  WITH
    TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (follower_id, following_id),
    FOREIGN KEY (follower_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (following_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS likes (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  post_id BIGINT,
  created_at TIMESTAMP(0)
  WITH
    TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, post_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
);

CREATE INDEX likes_user_post_idx ON likes (user_id, post_id);

CREATE INDEX idx_posts_id ON posts (id);
CREATE INDEX idx_posts_parent_post_id ON posts (parent_post_id);
CREATE INDEX idx_posts_created_at ON posts (created_at);
