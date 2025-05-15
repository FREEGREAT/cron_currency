package model

type Pair struct {
	ID    int    `json:"id"`
	Base  string `json:"base"`
	Quote string `json:"quote"`
}
