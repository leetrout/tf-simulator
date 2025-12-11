-- +goose Up
-- +goose StatementBegin
CREATE TABLE servers (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE servers;
-- +goose StatementEnd
