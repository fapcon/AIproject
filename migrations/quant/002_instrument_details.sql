-- +migrate Up
CREATE TABLE IF NOT EXISTS quant.instrument_details
(
    exchange        TEXT        NOT NULL,
    local_symbol    TEXT        NOT NULL,
    exchange_symbol TEXT        NOT NULL,
    price_precision INTEGER     NOT NULL,
    size_precision  INTEGER     NOT NULL,
    min_lot         DECIMAL     NOT NULL,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ          DEFAULT NULL,

    PRIMARY KEY (exchange, local_symbol, exchange_symbol)
);
-- +migrate Down
DROP TABLE IF EXISTS quant.instrument_details;