package handlers

import (
	"Go_project/offeringService/internal/app/api/requests"
	"Go_project/offeringService/internal/app/service"
	"Go_project/offeringService/pkg/models"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
)

var ErrorCloseReqBody = errors.New("failed to close request body")
var ErrorReadReqBody = errors.New("failed to read request body")

type OfferingHandler struct {
	Service *service.OfferingService
	Server  *http.Server
}

func NewHandler() *OfferingHandler {
	offerService := service.NewService()
	handler := OfferingHandler{Service: offerService}

	router := chi.NewRouter()
	router.Post("/offers", handler.CreateOffer)
	router.Get("/offers/{offerID}", handler.ParseOffer)

	handler.Server = &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	return &handler
}

func (handler *OfferingHandler) CreateOffer(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(ErrorReadReqBody, err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(ErrorCloseReqBody, err)
		}
	}(r.Body)

	var offerRequest requests.CreateOfferRequest
	err = json.Unmarshal(bytes, &offerRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	offer := models.Offer{From: offerRequest.From, To: offerRequest.To, ClientId: offerRequest.ClientId}
	offer.Price = *handler.Service.GetPrice(offer.From, offer.To)

	encodedJwt, err := handler.Service.EncodeJwt(&offer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	response := requests.CreateOfferResponse{OfferId: encodedJwt}
	res, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler *OfferingHandler) ParseOffer(w http.ResponseWriter, r *http.Request) {
	offerID := chi.URLParam(r, "offerID")

	decodedOffer, err := handler.Service.DecodeJwt(offerID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := requests.ParseOfferResponse{From: decodedOffer.From, To: decodedOffer.To, Price: decodedOffer.Price, ClientId: decodedOffer.ClientId}
	res, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
