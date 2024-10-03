-- name: SelectAllNameResolvers :many
SELECT *
FROM quant.name_resolvers
WHERE deleted_at IS NULL;


-- name: InsertNameResolvers :batchexec
INSERT INTO quant.name_resolvers (exchange, local_symbol, exchange_symbol)
VALUES ($1, $2, $3)
ON CONFLICT (exchange, local_symbol, exchange_symbol) DO UPDATE SET updated_at = NOW();
