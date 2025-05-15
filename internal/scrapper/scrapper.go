package scrapper

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/gocolly/colly"
)

type Rate struct {
	curr string
}

func ScrappUrl(url, domain string) (float64, error) {
	var resultt string
	var mu sync.Mutex
	c := colly.NewCollector(
		colly.AllowedDomains(domain),
	)

	c.OnHTML("span.ccOutputRslt", func(e *colly.HTMLElement) {
		text := e.DOM.Clone().Children().Remove().End().Text() // тільки основне число
		if text != "" {
			mu.Lock()
			resultt = text
			mu.Unlock()
		}
	})

	url = fmt.Sprintf("https://www.x-rates.com/calculator/?from=%s&to=%s&amount=1", "USD", "EUR")

	err := c.Visit(url)
	if err != nil {
		log.Fatalf("Error visiting base url: %v", err)
	}
	c.Wait()

	rate, err := strconv.ParseFloat(resultt, 64)
	if err != nil {
		return 0, fmt.Errorf("Error while parsing str to float %v", err)
	}

	return rate, nil

}
