package service

import (
	"context"
	"time"

	"scrapper.go/internal/model"
)

// type StorageService struct {
// 	currencyStorage storage.CurrencyStorage
// 	pairStorage     storage.PairStorage
// }

type StorageService interface {
	AddPair(ctx context.Context, base, quote string) error
	GetPairID(ctx context.Context, pair model.Pair) (int64, error)
	GetAllPairs(ctx context.Context) ([]model.Pair, error)
	SaveRate(ctx context.Context, pairID int64, rate float64, timestamp time.Time) error
	DeleteOldRates(ctx context.Context, pairID int64) error
	GetLatestRates(ctx context.Context, pairID int64) ([]model.Rate, error)
}

type ScrapService interface {
	FetchRate(base, quote string) (float64, error)
}
