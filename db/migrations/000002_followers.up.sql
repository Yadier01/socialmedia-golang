
CREATE TABLE IF NOT EXISTS followers (
    id BIGSERIAL PRIMARY KEY,
    follower_id BIGINT NOT NULL,
    following_id BIGINT NOT NULL, 
    UNIQUE(follower_id, following_id), 
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE
);

ALTER TABLE followers ADD CONSTRAINT unique_follow UNIQUE (follower_id, following_id);
