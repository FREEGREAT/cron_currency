package service

import (
	"context"
	"time"

	model "scrapper.go/internal/models"
	"scrapper.go/internal/storage"
)

type storageService struct {
	StoragePair     storage.PairStorage
	StorageCurrency storage.CurrencyStorage
}

// AddPair implements StorageService.
func (s *storageService) AddPair(ctx context.Context, base string, quote string) error {
	return s.StoragePair.AddPair(ctx, base, quote)
}

// DeleteOldRates implements StorageService.
func (s *storageService) DeleteOldRates(ctx context.Context, pairID int64) error {
	return s.StorageCurrency.DeleteOldRates(ctx, pairID)
}

// GetAllPairs implements StorageService.
func (s *storageService) GetAllPairs(ctx context.Context) ([]model.Pair, error) {
	return s.StoragePair.GetAllPairs(ctx)
}

// GetLatestRates implements StorageService.
func (s *storageService) GetLatestRates(ctx context.Context, pairID int64) ([]model.Rate, error) {
	return s.StorageCurrency.GetLatestRates(ctx, pairID)
}

// GetPairID implements StorageService.
func (s *storageService) GetPairID(ctx context.Context, pair model.Pair) (int64, error) {
	return s.StoragePair.GetPairID(ctx, pair)
}

// SaveRate implements StorageService.
func (s *storageService) SaveRate(ctx context.Context, pairID int64, rate float64, timestamp time.Time) error {
	return s.StorageCurrency.SaveRate(ctx, pairID, rate, timestamp)
}

func NewStorageService(currency storage.CurrencyStorage, pair storage.PairStorage) StorageService {
	return &storageService{
		StoragePair:     pair,
		StorageCurrency: currency,
	}
}
