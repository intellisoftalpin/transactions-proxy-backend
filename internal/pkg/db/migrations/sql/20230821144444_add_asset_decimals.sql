-- +goose Up
-- +goose StatementBegin
ALTER TABLE transactions ADD COLUMN asset_decimals VARCHAR NOT NULL DEFAULT '0';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE transactions DROP COLUMN asset_decimals;
-- +goose StatementEnd


