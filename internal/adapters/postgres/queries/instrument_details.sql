-- name: SelectAllInstrumentDetails :many
SELECT *
FROM quant.instrument_details
WHERE deleted_at IS NULL;


-- name: InsertInstrumentsDetails :batchexec
INSERT INTO quant.instrument_details (exchange, local_symbol, exchange_symbol, price_precision, size_precision, min_lot)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (exchange, local_symbol, exchange_symbol) DO UPDATE SET price_precision=EXCLUDED.price_precision,
                                                                    size_precision=EXCLUDED.size_precision,
                                                                    min_lot=EXCLUDED.min_lot,
                                                                    updated_at     = NOW();
