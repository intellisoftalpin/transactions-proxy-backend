-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS active_transactions (
    id                  SERIAL  PRIMARY KEY,
    user_id             INT    	NOT NULL,
    tx_id               INT     NOT NULL,
    step                VARCHAR NOT NULL, -- 'submit', 'updateStatus', 'send', 'onlyUpdateStatus'
    attempts            INT     NOT NULL DEFAULT 0,
    
    created_at          TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (tx_id) REFERENCES transactions(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS active_transactions;
-- +goose StatementEnd
