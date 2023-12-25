package models

type Price struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type Geolocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Offer struct {
	From     Geolocation
	To       Geolocation
	ClientId string
	Price    Price
}
