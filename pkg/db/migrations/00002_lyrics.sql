-- +goose Up
-- +goose StatementBegin
CREATE TABLE lyrics (
    id SERIAL PRIMARY KEY,
    song_id INT REFERENCES songs(id) ON DELETE CASCADE,
    verse_number INT,
    verse_lyrics TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE lyrics;
-- +goose StatementEnd
