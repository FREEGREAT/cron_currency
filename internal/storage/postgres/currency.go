package postgres

import (
	"context"
	"time"

	model "scrapper.go/internal/models"
	"scrapper.go/internal/storage"
	postgreConnect "scrapper.go/pkg/postgreSQL"
)

type currencyStorage struct {
	client postgreConnect.Client
}

func NewCurrencyRepository(client postgreConnect.Client) storage.CurrencyStorage {
	return &currencyStorage{
		client: client,
	}
}

// DeleteOldRates implements storage.CurrencyStorage.
func (c *currencyStorage) DeleteOldRates(ctx context.Context, pairID int64) error {
	q := `DELETE FROM currency_rates
        WHERE id NOT IN (
            SELECT id FROM currency_rates
            WHERE pair_id = $1
            ORDER BY timestamp DESC
            LIMIT 5
        ) AND pair_id = $1
	`

	_, err := c.client.Exec(ctx, q, pairID)
	if err != nil {
		return err
	}
	return nil
}

// GetLatestRates implements storage.CurrencyStorage.
func (c *currencyStorage) GetLatestRates(ctx context.Context, pairID int64) ([]model.Rate, error) {
	q := `
        SELECT rate, timestamp
		FROM currency_rates
		WHERE pair_id = $1
		ORDER BY timestamp DESC
		LIMIT 5
    `

	row, err := c.client.Query(ctx, q, pairID)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var rate []model.Rate
	for row.Next() {
		var rt model.Rate
		if err := row.Scan(&rt.Rate, &rt.Datetime); err != nil {
			return nil, err
		}
		rate = append(rate, rt)
	}
	return rate, nil
}

// SaveRate implements storage.CurrencyStorage.
func (c *currencyStorage) SaveRate(ctx context.Context, pairID int64, rate float64, timestamp time.Time) error {
	q := `INSERT INTO currency_rates (pair_id, rate) 
	VALUES ($1, $2)`

	_, err := c.client.Exec(ctx, q, pairID, rate)
	if err != nil {
		return err
	}
	return nil
}
