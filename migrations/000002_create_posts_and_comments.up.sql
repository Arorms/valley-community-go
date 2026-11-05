CREATE TABLE posts (
    id               BIGSERIAL PRIMARY KEY,
    author_id        UUID NOT NULL,
    title            VARCHAR(200) NOT NULL,
    markdown_content TEXT NOT NULL,
    html_content     TEXT NOT NULL DEFAULT '',
    like_count       INTEGER NOT NULL DEFAULT 0 CHECK (like_count >= 0),
    comment_count    INTEGER NOT NULL DEFAULT 0 CHECK (comment_count >= 0),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_deleted       BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_posts_created ON posts (created_at DESC) WHERE is_deleted = FALSE;
CREATE INDEX idx_posts_author ON posts (author_id, created_at DESC);

CREATE TABLE comments (
    id               BIGSERIAL PRIMARY KEY,
    post_id          BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    author_id        UUID NOT NULL,
    -- 两级结构：NULL 表示帖子的一级评论，否则指向一级评论
    parent_id        BIGINT REFERENCES comments(id) ON DELETE CASCADE,
    -- 回复的回复仍归到同一线程，用该字段记录被回复的人
    reply_to_user_id UUID,
    content_markdown TEXT NOT NULL,
    content_html     TEXT NOT NULL DEFAULT '',
    like_count       INTEGER NOT NULL DEFAULT 0 CHECK (like_count >= 0),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_deleted       BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT chk_comments_parent_not_self CHECK (parent_id IS NULL OR parent_id <> id)
);

CREATE INDEX idx_comments_post ON comments (post_id, created_at) WHERE is_deleted = FALSE;
CREATE INDEX idx_comments_parent ON comments (parent_id) WHERE parent_id IS NOT NULL;
CREATE INDEX idx_comments_author ON comments (author_id, created_at DESC);