package api

import (
	client "Go_project/clientTests"
	"Go_project/internal/app/api/requests"
	"Go_project/pkg/models"
	"encoding/json"
	"fmt"
)

func TestRequests(client *client.Client) {
	requestCreate := requests.CreateOfferRequest{
		From:     models.Geolocation{Lat: 55, Lng: 33},
		To:       models.Geolocation{Lat: 59, Lng: 30},
		ClientId: "1",
	}

	result1, err := client.GetOffer(requestCreate)
	if err != nil {
		fmt.Println(err)
	} else {
		bytes, _ := json.Marshal(result1)
		fmt.Println(string(bytes))
	}

	requestParse := requests.ParseOfferRequest{OfferId: result1.OfferId}

	result2, err := client.GetParsedOffer(requestParse)
	if err != nil {
		fmt.Println(err)
	} else {
		bytes, _ := json.Marshal(result2)
		fmt.Println(string(bytes))
	}
}
