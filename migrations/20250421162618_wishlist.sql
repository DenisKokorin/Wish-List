-- +goose Up
CREATE TABLE wishlist (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    title VARCHAR NOT NULL,
    is_private BOOLEAN DEFAULT false
);

CREATE TABLE item (
    id BIGSERIAL PRIMARY KEY,
    wishlist_id BIGINT REFERENCES wishlist (id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255)
);

-- +goose Down
DROP TABLE wishlist;
DROP TABLE item;

