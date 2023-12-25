package client

import (
	"Go_project/internal/app/api/requests"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	url string
}

func NewClient(url string) *Client {
	return &Client{url}
}

func (client *Client) GetOffer(input requests.CreateOfferRequest) (requests.CreateOfferResponse, error) {
	inputJson, err := json.Marshal(input)
	if err != nil {
		return requests.CreateOfferResponse{}, err
	}

	resp, err := http.Post(client.url+"/offers", "application/json", bytes.NewBuffer(inputJson))
	if err != nil {
		return requests.CreateOfferResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return requests.CreateOfferResponse{}, err
	}

	var result requests.CreateOfferResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return requests.CreateOfferResponse{}, fmt.Errorf("Status: %v Error: %v", resp.Status, string(body))
	}

	return result, nil
}

func (client *Client) GetParsedOffer(input requests.ParseOfferRequest) (requests.ParseOfferResponse, error) {
	resp, err := http.Get(client.url + "/offers/" + input.OfferId)
	if err != nil {
		return requests.ParseOfferResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return requests.ParseOfferResponse{}, err
	}

	var result requests.ParseOfferResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return requests.ParseOfferResponse{}, fmt.Errorf("Status: %v Error: %v", resp.Status, string(body))
	}

	return result, nil
}
