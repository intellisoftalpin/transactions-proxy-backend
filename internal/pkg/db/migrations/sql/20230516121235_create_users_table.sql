-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
	id          	SERIAL  PRIMARY KEY,
	user_hash       VARCHAR NOT NULL,
	user_runtime    INT2    NOT NULL,
	created_at     	TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
