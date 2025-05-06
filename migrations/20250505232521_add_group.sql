-- +goose Up
CREATE TABLE groups (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT REFERENCES users (id) ON DELETE CASCADE,
    title VARCHAR NOT NULL
);

CREATE TABLE group_member (
    group_id BIGINT REFERENCES groups (id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE groups;
DROP TABLE group_member;