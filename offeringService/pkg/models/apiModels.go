package models

type Geolocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Price struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type Offer struct {
	From     Geolocation
	To       Geolocation
	ClientId string `json:"client_id"`
	Price    Price
}
