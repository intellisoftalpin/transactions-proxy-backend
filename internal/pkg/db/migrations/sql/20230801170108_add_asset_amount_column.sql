-- +goose Up
-- +goose StatementBegin
ALTER TABLE transactions ADD COLUMN asset_amount VARCHAR NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE transactions DROP COLUMN asset_amount;
-- +goose StatementEnd
