-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transactions (
    id                  SERIAL  PRIMARY KEY,
    user_id             INT    	NOT NULL,
    tx_type             VARCHAR NOT NULL,
    tx_status           VARCHAR NOT NULL DEFAULT 'draft',
    tx_hash             VARCHAR NOT NULL DEFAULT '',
    tx_cbor             VARCHAR NOT NULL DEFAULT '',
    addr_to             VARCHAR NOT NULL DEFAULT '',
    transfer_amount     VARCHAR NOT NULL DEFAULT '',
    policy_id           VARCHAR NOT NULL DEFAULT '',
    asset_id            VARCHAR NOT NULL DEFAULT '',
    created_at          TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd
