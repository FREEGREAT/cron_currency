package service

import (
	"fmt"

	scraper "scrapper.go/internal/scrapper"
)

const (
	domain  = "www.x-rates.com"
	zeroVal = 0
)

func FetchRate(base, quote string) (float64, error) {
	url := fmt.Sprintf("https://www.x-rates.com/calculator/?from=%s&to=%s&amount=1", base, quote)

	fmt.Printf("Base %s, Quote %s", base, quote)

	rate, err := scraper.ScrappUrl(url, domain)
	if err != nil {
		return zeroVal, err
	}

	fmt.Println(rate)

	return rate, nil
}
