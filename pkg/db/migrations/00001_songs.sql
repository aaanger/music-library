-- +goose Up
-- +goose StatementBegin
CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    song VARCHAR(255) NOT NULL,
    artist VARCHAR(255) NOT NULL,
    release_date VARCHAR(255),
    link TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE songs;
-- +goose StatementEnd
