CREATE TABLE notifications (
    id           BIGSERIAL PRIMARY KEY,
    recipient_id UUID NOT NULL,
    actor_id     UUID NOT NULL,
    type         VARCHAR(32) NOT NULL,
    post_id      BIGINT REFERENCES posts(id) ON DELETE CASCADE,
    comment_id   BIGINT REFERENCES comments(id) ON DELETE CASCADE,
    read_at      TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_notifications_type CHECK (
        type IN ('comment', 'reply', 'mention', 'post_like', 'comment_like')
    ),
    CONSTRAINT chk_notifications_target CHECK (post_id IS NOT NULL OR comment_id IS NOT NULL)
);

CREATE INDEX idx_notifications_recipient ON notifications (recipient_id, created_at DESC);
CREATE INDEX idx_notifications_unread ON notifications (recipient_id) WHERE read_at IS NULL;