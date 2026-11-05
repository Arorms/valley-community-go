CREATE TABLE post_likes (
    post_id    BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (post_id, user_id)
);

CREATE INDEX idx_post_likes_user ON post_likes (user_id, created_at DESC);

CREATE TABLE comment_likes (
    comment_id BIGINT NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (comment_id, user_id)
);

CREATE INDEX idx_comment_likes_user ON comment_likes (user_id, created_at DESC);

-- 每个用户一个隐式收藏夹：复合主键即可，无需额外的收藏夹表
CREATE TABLE favorites (
    user_id UUID NOT NULL,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, post_id)
);

CREATE INDEX idx_favorites_post ON favorites (post_id);