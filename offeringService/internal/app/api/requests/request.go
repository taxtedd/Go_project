package requests

import (
	"Go_project/offeringService/pkg/models"
)

type CreateOfferRequest struct {
	From     models.Geolocation `json:"from"`
	To       models.Geolocation `json:"to"`
	ClientId string             `json:"client_Id"`
}

type ParseOfferRequest struct {
	OfferId string `json:"offer_Id"`
}

type ParseOfferResponse struct {
	From     models.Geolocation `json:"from"`
	To       models.Geolocation `json:"to"`
	ClientId string             `json:"client_Id"`
	Price    models.Price       `json:"price"`
}

type CreateOfferResponse struct {
	OfferId string `json:"offer_Id"`
}
