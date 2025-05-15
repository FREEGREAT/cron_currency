package model

import "time"

type Rate struct {
	Rate     float64   `json:"rate"`
	Datetime time.Time `json:"timestamp"`
}
