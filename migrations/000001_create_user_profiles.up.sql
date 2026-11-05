CREATE TABLE user_profiles (
    user_id    UUID PRIMARY KEY,
    username   VARCHAR(64) NOT NULL,
    nickname   VARCHAR(64) NOT NULL DEFAULT '',
    avatar_url TEXT NOT NULL DEFAULT '',
    gender     VARCHAR(16) NOT NULL DEFAULT '',
    birthday   DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX uq_user_profiles_username ON user_profiles (lower(username));