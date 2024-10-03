-- +migrate Up
CREATE SCHEMA IF NOT EXISTS quant;
CREATE TABLE IF NOT EXISTS quant.name_resolvers
(
    exchange        TEXT        NOT NULL,
    local_symbol    TEXT        NOT NULL,
    exchange_symbol TEXT        NOT NULL,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ          DEFAULT NULL,

    PRIMARY KEY (exchange, local_symbol, exchange_symbol)
);
-- +migrate Down
DROP TABLE IF EXISTS quant.name_resolvers;
-- DROP SCHEMA quant CASCADE;