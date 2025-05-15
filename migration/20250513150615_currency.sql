-- +goose Up
-- +goose StatementBegin
CREATE TABLE subscribed_pairs (
    id SERIAL PRIMARY KEY,
    base_currency VARCHAR(3),
    quote_currency VARCHAR(3),
    UNIQUE(base_currency, quote_currency)
);

CREATE TABLE currency_rates (
    id SERIAL PRIMARY KEY,
    pair_id INT REFERENCES subscribed_pairs(id) ON DELETE CASCADE,
    rate NUMERIC(10, 4),
    timestamp TIMESTAMPTZ DEFAULT now()
);



-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE exchange_rates;
DROP TABLE currency_pairs;
-- +goose StatementEnd